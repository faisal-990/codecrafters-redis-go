package command

import (
	"github.com/codecrafters-io/redis-starter-go/db"
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

func SetHandler(cmd *resp.Command) resp.Resp {
	err := db.Instance.Set(cmd.Args[0], cmd.Args[1])
	if err != nil {
		return resp.RespError("ERR - Failed to set the key value pair")
	}
	// when V field of simple string response is passed as nil
	// the response that simple string sends in +OK, if some string is passed
	// then that string is sent as response
	return resp.SimpleString{V: nil}
}

func GetHandler(cmd *resp.Command) resp.Resp {
	value, err := db.Instance.Get(cmd.Args[0])
	if err != nil {
		return resp.RespError("ERR- failed to get the value using the key")
	}
	if value == "" {
		// in case the key doesn't exist in the memory return a null bulk string
		return resp.BulkString{V: nil}
	}
	return resp.BulkString{V: &(value)}
}
