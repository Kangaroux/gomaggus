package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/internal/worldd"
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
	server := worldd.NewServer(db, worldd.DefaultPort)
	server.Start()
}
