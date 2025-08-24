package command

import (
	"strconv"
	"strings"
	"time"

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
	key := cmd.Args[0]
	value := cmd.Args[1]
	var timeout time.Duration
	// now parsing the flags which are optional
	if len(cmd.Args) > 2 {
		flag := strings.ToUpper(cmd.Args[2])
		if flag == "PX" {
			ms, err := strconv.Atoi(cmd.Args[3])
			if err != nil {
				return resp.RespError("ERR - invalid time format")
			}
			timeout = time.Duration(ms) * time.Millisecond

		}

	}
	err := db.Instance.Set(key, value, timeout)
	if err != nil {
		return resp.RespError("Err - not able to store in databse")
	}
	// when V field of simple string response is passed as nil
	// the response that simple string sends in +OK, if some string is passed
	// then that string is sent as response
	return resp.SimpleString{V: nil}
}

func GetHandler(cmd *resp.Command) resp.Resp {
	value, err := db.Instance.Get(cmd.Args[0])

	if value == "" {
		// in case the key doesn't exist in the memory return a null bulk string
		return resp.BulkString{V: nil}
	}
	if err != nil {
		return resp.RespError("ERR- failed to get the value using the key")
	}
	return resp.BulkString{V: &(value)}
}

func RpushHandler(cmd *resp.Command) resp.Resp {
	key := cmd.Args[0]
	values := cmd.Args[1:]

	length, err := db.Instance.Rpush(key, values)
	if err != nil {
		if strings.HasPrefix(err.Error(), "WRONGTYPE") {
			return resp.RespError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return resp.RespError(err.Error())
	}
	return resp.Integer(strconv.Itoa(length))
}

func LrangeHandler(cmd *resp.Command) resp.Resp {
	key := cmd.Args[0]
	start, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		return resp.RespError("ERR - value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		return resp.RespError("ERR - value is not an integer or out of range")
	}

	values, err := db.Instance.Lrange(key, start, stop)
	if len(values) == 0 {
		// no elements found for what is being queried ,return a simple empty array
		return resp.RespArray{V: []resp.Resp{}}
	}
	//if err != nil {
	//return resp.RespError(err.Error())
	//}

	// convert []string -> []Resp (BulkString)
	arr := make([]resp.Resp, len(values))
	for i, v := range values {
		val := v // capture for pointer
		arr[i] = resp.BulkString{V: &val}
	}

	return resp.RespArray{V: arr}
}
