package main

import (
	"bufio"
	"crypto"
	"fmt"
	"os"
	"strings"

	_ "crypto/sha256"

	srp "github.com/arag0re/go-apple-srp6"
)

const (
	MinAccountNameLen = 3
	MaxAccountNameLen = 12
	MinPasswordLen    = 8
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("username: ")

	accountName, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	accountName = strings.ToUpper(strings.TrimSpace(accountName))

	if len(accountName) < MinAccountNameLen || len(accountName) > MaxAccountNameLen {
		fmt.Printf("error: username must be between %d and %d characters\n", MinAccountNameLen, MaxAccountNameLen)
		os.Exit(1)
	}

	fmt.Print("password: ")

	password, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	password = strings.TrimSpace(password)

	if len(password) < MinPasswordLen {
		fmt.Printf("error: password must be at least %d characters\n", MinPasswordLen)
		os.Exit(1)
	}

	s, err := srp.NewWithHash(crypto.SHA256, 4096)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v, err := s.Verifier([]byte(accountName), []byte(password))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, verif := v.Encode()
	parts := strings.Split(verif, ":")
	fmt.Println("s:", parts[5])
	fmt.Println("v:", parts[6])
}
