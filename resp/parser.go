package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

// ParseCommand parses a Redis array command
func ParseCommand(r *bufio.Reader) (*Command, error) {
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
