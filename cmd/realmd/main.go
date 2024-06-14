package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/realmd/server"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect(
		"postgres",
		"postgres://gomaggus:password@localhost:5432/gomaggus?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	server := server.New(db, server.DefaultListenAddr)
	server.Start()
}
