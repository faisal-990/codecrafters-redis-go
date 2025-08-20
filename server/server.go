package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/command"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		cmd, err := resp.ParseCommand(r)
		if err != nil {
			if err == io.EOF {
				return
			}
			w.WriteString(resp.RespError(fmt.Sprintf("ERR %v", err)).ToRESP())
			w.Flush()
			continue
		}

		w.WriteString(command.Dispatch(cmd).ToRESP())
		w.Flush()
	}
}

func StartServer(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
