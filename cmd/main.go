package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/server"
)

func main() {
	addr := ":6379"
	fmt.Println("listening on", addr)
	if err := server.StartServer(addr); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
