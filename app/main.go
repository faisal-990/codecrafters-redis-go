package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// ---------- core types ----------
type Command struct {
	Name string
	Args []string
}

type Resp interface {
	ToRESP() string
}

// BulkString Resp implementation (nil -> $-1)
type BulkString struct{ V *string }

func (b BulkString) ToRESP() string {
	if b.V == nil {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(*b.V), *b.V)
}

// RespError for protocol / unknown-command errors
type RespError string

func (e RespError) ToRESP() string { return "-" + string(e) + "\r\n" }

// Handler signature
type Handler func(cmd *Command) Resp

// ---------- command registry (extend here) ----------
var cmdTable = map[string]Handler{
	"ECHO": echoHandler,
	"PING": pingHandler,
}

func pingHandler(cmd *Command) Resp {
	str := "PONG"
	return BulkString{V: &str}
}

func echoHandler(cmd *Command) Resp {
	if len(cmd.Args) < 1 {
		return RespError("ERR wrong number of arguments for 'echo'")
	}
	s := cmd.Args[0]
	return BulkString{V: &s}
}

func dispatch(cmd *Command) Resp {
	if h, ok := cmdTable[strings.ToUpper(cmd.Name)]; ok {
		return h(cmd)
	}
	return RespError(fmt.Sprintf("ERR unknown command '%s'", cmd.Name))
}

// ---------- RESP parsing helpers ----------
func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")
	return s, nil
}

func readCRLF(r *bufio.Reader) error {
	b := make([]byte, 2)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	if b[0] != '\r' || b[1] != '\n' {
		return fmt.Errorf("expected CRLF, got %q", string(b))
	}
	return nil
}

func parseBulkString(r *bufio.Reader) (string, error) {
	// expect $
	b, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	if b != '$' {
		return "", fmt.Errorf("expected '$', got %q", b)
	}

	lenLine, err := readLine(r)
	if err != nil {
		return "", err
	}
	n, err := strconv.Atoi(lenLine)
	if err != nil {
		return "", fmt.Errorf("invalid bulk length %q", lenLine)
	}

	// null bulk
	if n < 0 {
		return "", nil
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}

	if err := readCRLF(r); err != nil {
		return "", err
	}
	return string(buf), nil
}

func parseCommand(r *bufio.Reader) (*Command, error) {
	// expect '*'
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != '*' {
		return nil, fmt.Errorf("expected '*', got %q", b)
	}

	lenLine, err := readLine(r)
	if err != nil {
		return nil, err
	}
	count, err := strconv.Atoi(lenLine)
	if err != nil {
		return nil, fmt.Errorf("invalid array length %q", lenLine)
	}
	if count <= 0 {
		return nil, fmt.Errorf("array length must be > 0")
	}

	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		arg, err := parseBulkString(r)
		if err != nil {
			return nil, err
		}
		parts = append(parts, arg)
	}
	return &Command{Name: strings.ToUpper(parts[0]), Args: parts[1:]}, nil
}

// ---------- connection loop (handles pipelining) ----------
func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		// parse first command (blocks until a complete one is available)
		cmd, err := parseCommand(r)
		if err != nil {
			if err == io.EOF {
				return
			}
			w.WriteString(RespError(fmt.Sprintf("ERR %v", err)).ToRESP())
			w.Flush()
			continue
		}

		// dispatch & write reply
		w.WriteString(dispatch(cmd).ToRESP())

		// handle any other complete commands already buffered (pipelining)
		//for r.Buffered() > 0 {
		//cmd2, err := parseCommand(r)
		//if err != nil {
		//// protocol error for the pipelined command -> return an error reply and stop draining
		//w.WriteString(RespError(fmt.Sprintf("ERR %v", err)).ToRESP())
		//break
		//}
		//w.WriteString(dispatch(cmd2).ToRESP())
		//}
		//
		w.Flush()
	}
}

// ---------- main ----------package main

// StartServer runs your TCP server (can be called by main or tests)
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
		go handleConnection(conn) // your existing connection handler
	}
}

func main() {
	addr := ":6379"
	fmt.Println("listening on", addr)
	if err := StartServer(addr); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
