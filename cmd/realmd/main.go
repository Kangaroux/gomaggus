package main

import "github.com/kangaroux/gomaggus/realmd"

func main() {
	server := realmd.NewServer(realmd.DefaultPort)
	server.Start()
}
