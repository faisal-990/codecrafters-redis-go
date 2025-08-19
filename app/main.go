package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)
fsdafa
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	fmt.Print(reader.Size())
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		fmt.Println("Failed to bind to port 6380")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
