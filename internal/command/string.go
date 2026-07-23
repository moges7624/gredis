package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/moges7624/gredis/internal/resp"
	"github.com/moges7624/gredis/internal/store"
)

func handlePing(s *store.Store, args []string) []byte {
	if len(args) > 1 {
		return resp.EncodeError("wrong number of arguments for 'ping' command")
	}

	if len(args) == 1 {
		return resp.EncodeBulkString(args[0])
	}

	return resp.EncodeSimpleString("PONG")
}

func handleEcho(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeError("wrong number of arguments for 'echo' command")
	}

	return resp.EncodeBulkString(args[0])
}

func handleInfo(s *store.Store, args []string) []byte {
	info := "# Server\r\nredis_version:7.2.0\r\ntcp_port:6379\r\n"
	return resp.EncodeBulkString(info)
}

func handleGet(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeError("wrong number of arguments for 'get' command")
	}

	val, exists, err := s.Get(args[0])
	if err != nil {
		return resp.EncodeError(err.Error())
	}

	if !exists {
		return resp.EncodeNullString()
	}

	return resp.EncodeBulkString(val)
}

func handleSet(s *store.Store, args []string) []byte {
	if len(args) < 2 {
		return resp.EncodeError("wrong number of arguments for 'set' command")
	}

	key, val := args[0], args[1]
	var opts store.SetOptions
	for i := 2; i < len(args); i++ {
		switch strings.ToUpper(args[i]) {
		case "EX":
			if len(args) <= i+1 {
				return resp.EncodeError("syntax error")
			}
			n, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil || n <= 0 {
				return resp.EncodeError("invalid expire time in 'set' command")
			}

			opts.TTL = time.Duration(n) * time.Second
			i++
		}
	}

	s.Set(key, val, opts)

	return resp.EncodeSimpleString("OK")
}

func handleDel(s *store.Store, args []string) []byte {
	if len(args) < 1 {
		return resp.EncodeError("wrong number of arguments for 'del' command")
	}

	count := s.Delete(args)
	return resp.EncodeInteger(int64(count))
}

func handleIncr(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeError("wrong number of arguments for 'incr' command")
	}

	val, exists, err := s.IncrBy(args[0], 1)
	if exists && err != nil {
		return resp.EncodeError(err.Error())
	}

	if !exists {
		s.Set(args[0], strconv.Itoa(1), store.SetOptions{})
		return resp.EncodeInteger(1)
	}

	return resp.EncodeInteger(val)
}

func handleDecr(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeError("wrong number of arguments for 'decr' command")
	}

	val, exists, err := s.IncrBy(args[0], -1)
	if exists && err != nil {
		return resp.EncodeError(err.Error())
	}

	if !exists {
		s.Set(args[0], strconv.Itoa(-1), store.SetOptions{})
		return resp.EncodeInteger(-1)
	}

	return resp.EncodeInteger(val)
}
