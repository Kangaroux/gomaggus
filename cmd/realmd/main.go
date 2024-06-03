package main

import "github.com/kangaroux/gomaggus/internal/realmd"

func main() {
	server := realmd.NewServer(realmd.DefaultPort)
	server.Start()
}
