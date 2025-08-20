package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestEchoCommands(t *testing.T) {
	tests := []struct {
		name      string
		input     string // RESP command
		wantLines []string
	}{
		{
			name:      "simple echo",
			input:     "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
			wantLines: []string{"$5\r\n", "hello\r\n"},
		},
		{
			name:      "empty echo",
			input:     "*2\r\n$4\r\nECHO\r\n$0\r\n\r\n",
			wantLines: []string{"$0\r\n", "\r\n"},
		},
		{
			name:      "missing argument",
			input:     "*1\r\n$4\r\nECHO\r\n",
			wantLines: []string{"-ERR wrong number of arguments for 'echo'\r\n"},
		},
	}

	addr := "127.0.0.1:6381" // use a separate port for tests
	go StartServer(addr)
	time.Sleep(100 * time.Millisecond)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			conn.Write([]byte(tt.input))
			reader := bufio.NewReader(conn)

			for _, wantLine := range tt.wantLines {
				gotLine, _ := reader.ReadString('\n')
				if gotLine != wantLine {
					t.Errorf("expected %q, got %q", wantLine, gotLine)
				}
			}
		})
	}
}
