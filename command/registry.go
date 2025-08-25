package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/resp"
)

type Handler func(cmd *resp.Command) resp.Resp

var cmdTable = map[string]Handler{
	"ECHO":   EchoHandler,
	"PING":   PingHandler,
	"SET":    SetHandler,
	"GET":    GetHandler,
	"RPUSH":  RpushHandler,
	"LRANGE": LrangeHandler,
	"LPUSH":  LpushHandler,
	"LLEN":   LlenHandler,
	"LPOP":   LpopHandler,
}

func Dispatch(cmd *resp.Command) resp.Resp {
	if h, ok := cmdTable[strings.ToUpper(cmd.Name)]; ok {
		return h(cmd)
	}
	return resp.RespError(fmt.Sprintf("ERR unknown command '%s'", cmd.Name))
}
