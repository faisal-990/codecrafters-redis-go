package command

import (
	"github.com/codecrafters-io/redis-starter-go/resp"
)

func EchoHandler(cmd *resp.Command) resp.Resp {
	if len(cmd.Args) < 1 {
		return resp.RespError("ERR wrong number of arguments for 'echo'")
	}
	s := cmd.Args[0]
	return resp.BulkString{V: &s}
}

func PingHandler(cmd *resp.Command) resp.Resp {
	str := "PONG"
	return resp.SimpleString{V: &str}
}
