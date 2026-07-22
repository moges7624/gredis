package command

import (
	"fmt"
	"strconv"

	"github.com/moges7624/gredis/internal/store"
)

func handlePing(s *store.Store, args []string) []byte {
	return []byte("+PONG\r\n")
}

func handleInfo(s *store.Store, args []string) []byte {
	CRLF := "\r\n"
	resp := "# Server\r\nredis_version:7.2.0\r\ntcp_port:6379\r\n"
	resp = ("$" + strconv.Itoa(len(resp)) + CRLF + string(resp) + CRLF)
	return []byte(resp)
}

func handleGet(s *store.Store, args []string) []byte {
	val, exists := s.Get(args[0])
	if !exists {
		return []byte("$-1\r\n")
	}

	return fmt.Appendf(nil, "+%s\r\n", val)
}

func handleSet(s *store.Store, args []string) []byte {
	if len(args) < 2 {
		return []byte("-ERR wrong number of arguments for 'set' command")
	}

	key, val := args[0], args[1]
	s.Set(key, val)

	return []byte("+OK\r\n")
}
