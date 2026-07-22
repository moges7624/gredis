package command

import (
	"fmt"

	"github.com/moges7624/gredis/internal/resp"
	"github.com/moges7624/gredis/internal/store"
)

func handlePing(s *store.Store, args []string) []byte {
	return []byte("+PONG\r\n")
}

func handleInfo(s *store.Store, args []string) []byte {
	info := "# Server\r\nredis_version:7.2.0\r\ntcp_port:6379\r\n"
	return resp.EncodeBulkString(info)
}

func handleGet(s *store.Store, args []string) []byte {
	val, exists := s.Get(args[0])
	if !exists {
		return resp.EncodeBulkString("-1")
	}

	return fmt.Appendf(nil, "+%s\r\n", val)
}

func handleSet(s *store.Store, args []string) []byte {
	if len(args) < 2 {
		return resp.EncodeError("wrong number of arguments for 'set' command")
	}

	key, val := args[0], args[1]
	s.Set(key, val)

	return resp.EncodeBulkString("OK")
}
