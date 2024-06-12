package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/internal/realmd"
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
	server := realmd.NewServer(db, realmd.DefaultListenAddr)
	server.Start()
}
