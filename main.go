package main

import (
	"fmt"
	"web-ssh-server/bootstrap"
	"web-ssh-server/config"
)

func main() {
	// 1. Parse config.
	config.ParseConfig()

	// 2. Startup server.
	bootstrap.Run()

	fmt.Printf("Started.")
}
