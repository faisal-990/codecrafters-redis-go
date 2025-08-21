package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/db"
	"github.com/codecrafters-io/redis-starter-go/server"
)

func main() {
	addr := ":6379"
	fmt.Println("listening on", addr)

	// initializing the in memory databse , this db is passed by reference throughout the application
	// to use db.Instance gives you access to its methods
	db.Init()

	if err := server.StartServer(addr); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
