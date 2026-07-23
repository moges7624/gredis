package command

import (
	"github.com/moges7624/gredis/internal/resp"
	"github.com/moges7624/gredis/internal/store"
)

func handleTTL(s *store.Store, args []string) []byte {
	if len(args) > 1 {
		return resp.EncodeError("wrong number of arguments for 'ttl' command")
	}

	ttl, exists, hasTTL := s.TTL(args[0])
	if !exists {
		return resp.EncodeInteger(-2)
	}

	if !hasTTL {
		return resp.EncodeInteger(-1)
	}

	secs := int64(ttl.Seconds())
	if secs == 0 && ttl > 0 {
		secs = 1
	}

	return resp.EncodeInteger(secs)
}

func handleExists(s *store.Store, args []string) []byte {
	if len(args) == 0 {
		return resp.EncodeError("wrong number of arguments for 'exists' command")
	}

	cnt := s.Exists(args)
	return resp.EncodeInteger(int64(cnt))
}

func handleType(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeError("wrong number of arguments for 'type' command")
	}

	t := s.Type(args[0])
	return resp.EncodeSimpleString(t)
}
