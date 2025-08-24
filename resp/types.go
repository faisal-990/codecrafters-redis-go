package resp

import (
	"fmt"
	"strings"
)

// Command represents a parsed Redis command
type Command struct {
	Name string
	Args []string
}

// Resp interface for all response types
// ToRESP method prepares the response in the resp protocol format
type Resp interface {
	ToRESP() string
}

// BulkString implementation
type BulkString struct{ V *string }

func (b BulkString) ToRESP() string {
	if b.V == nil {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(*b.V), *b.V)
}

// SimpleString implementation
type SimpleString struct{ V *string }

func (s SimpleString) ToRESP() string {
	if s.V == nil {
		return "+OK\r\n"
	}
	return fmt.Sprintf("+%s\r\n", *s.V)
}

// RespError for errors
type RespError string

func (e RespError) ToRESP() string { return "-" + string(e) + "\r\n" }

type Integer string

func (i Integer) ToRESP() string {
	return ":" + string(i) + "\r\n"
}

type RespArray struct {
	V []Resp
}

func (r RespArray) ToRESP() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(r.V)))
	for _, item := range r.V {
		sb.WriteString(item.ToRESP())
	}
	return sb.String()
}
