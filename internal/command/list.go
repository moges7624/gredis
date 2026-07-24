package command

import (
	"github.com/moges7624/gredis/internal/resp"
	"github.com/moges7624/gredis/internal/store"
)

func handleLLen(s *store.Store, args []string) []byte {
	if len(args) != 1 {
		return resp.EncodeWrongArgumentNumber("llen")
	}

	currLen, err := s.LLen(args[0])
	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(int64(currLen))
}

func handleRPush(s *store.Store, args []string) []byte {
	if len(args) < 2 {
		return resp.EncodeWrongArgumentNumber("rpush")
	}

	key, vals := args[0], args[1:]

	count, err := s.RPush(key, vals...)
	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(int64(count))
}
