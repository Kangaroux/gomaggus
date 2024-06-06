package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/internal/models"
	_ "github.com/lib/pq"
)

func usage() {
	fmt.Println("usage:", os.Args[0], "<command>")
	fmt.Println()
	fmt.Println("commands")
	fmt.Println()
	fmt.Println("    add, a           Add a new account")
	fmt.Println("    password, p      Change an existing account's password")
	fmt.Println()
}

func addUsage() {
	fmt.Println("usage:", os.Args[1], "<username> <password> <email>")
}

func passwordUsage() {
	fmt.Println("usage:", os.Args[1], "<username> <password>")
}

func main() {
	db, err := sqlx.Connect(
		"postgres",
		"postgres://gomaggus:password@localhost:5432/gomaggus?sslmode=disable",
	)
	if err != nil {
		fmt.Println("failed to connect to db:", err)
		os.Exit(1)
	}

	accountsDb := models.NewDbAccountService(db)

	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add", "a":
		args := os.Args[2:]

		if len(args) != 4 {
			fmt.Println("error: expected 3 arguments")
			addUsage()
			os.Exit(1)
		}

		username := strings.TrimSpace(args[0])
		password := args[1]
		email := strings.TrimSpace(args[2])

		if len(username) < 3 || len(username) > 16 {
			fmt.Println("error: username must be between 3-16 characters")
			addUsage()
			os.Exit(1)
		} else if len(password) < 6 || len(password) > 16 {
			fmt.Println("error: password must be between 6-16 characters")
			addUsage()
			os.Exit(1)
		}

		existingAccount, err := accountsDb.Get(&models.AccountGetParams{
			Email:    email,
			Username: username,
		})
		if err != nil {
			fmt.Println("failed to get account:", err)
			os.Exit(1)
		} else if existingAccount != nil {
			fmt.Println("error: username or email is already taken")
			os.Exit(1)
		}

		account := &models.Account{
			Username: username,
			Email:    email,
		}
		account.SetUsernamePassword(username, password)

		if err := accountsDb.Create(account); err != nil {
			fmt.Println("failed to create account:", err)
			os.Exit(1)
		}

		fmt.Println("success")
		fmt.Println("account id:", account.Id)

	case "password", "p":
		args := os.Args[2:]

		if len(args) != 2 {
			fmt.Println("error: expected 2 arguments")
			passwordUsage()
			os.Exit(1)
		}

		username := strings.TrimSpace(args[0])
		password := args[1]

		if len(password) < 6 || len(password) > 16 {
			fmt.Println("error: password must be between 6-16 characters")
			passwordUsage()
			os.Exit(1)
		}

		account, err := accountsDb.Get(&models.AccountGetParams{Username: username})
		if err != nil {
			fmt.Println("failed to get account:", err)
			os.Exit(1)
		} else if account == nil {
			fmt.Println("error: no account with that username exists")
			os.Exit(1)
		}

		account.SetUsernamePassword(username, password)
		if _, err := accountsDb.Update(account); err != nil {
			fmt.Println("failed to update account:", err)
			os.Exit(1)
		}

		fmt.Println("success")
	}
}
