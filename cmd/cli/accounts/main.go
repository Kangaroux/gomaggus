package main

import (
	"fmt"
	"os"
	"strconv"
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
	fmt.Println("\tadd, a")
	fmt.Println()
}

func addUsage() {
	fmt.Println("usage:", os.Args[1], "<username> <password> <email> <realmId>")
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
	realmsDb := models.NewDbRealmService(db)

	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add", "a":
		args := os.Args[2:]

		if len(args) != 4 {
			fmt.Println("error: expected 4 arguments")
			addUsage()
			os.Exit(1)
		}

		username := strings.TrimSpace(args[0])
		password := args[1]
		email := strings.TrimSpace(args[2])
		realmId, err := strconv.Atoi(args[3])

		if len(username) < 3 || len(username) > 16 {
			fmt.Println("error: username must be between 3-16 characters")
			addUsage()
			os.Exit(1)
		} else if len(password) < 6 || len(password) > 16 {
			fmt.Println("error: password must be between 6-16 characters")
			addUsage()
			os.Exit(1)
		} else if err != nil {
			fmt.Println("failed to parse realm id:", err)
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

		realm, err := realmsDb.Get(uint32(realmId))
		if err != nil {
			fmt.Println("failed to get realm:", err)
			os.Exit(1)
		} else if realm == nil {
			fmt.Printf("error: no realm with id %d exists\n", realmId)
			os.Exit(1)
		}

		account := &models.Account{
			Username: username,
			Email:    email,
			RealmId:  uint32(realmId),
		}
		account.SetUsernamePassword(username, password)

		if err := accountsDb.Create(account); err != nil {
			fmt.Println("failed to create account:", err)
			os.Exit(1)
		}

		fmt.Println("success")
		fmt.Println("account id:", account.Id)
	}
}
