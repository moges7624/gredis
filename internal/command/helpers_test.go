package command

import (
	"strconv"
	"strings"
	"testing"

	"github.com/moges7624/gredis/internal/store"
)

func newDispatcher(t *testing.T) *Dispatcher {
	t.Helper()
	db := store.NewStore()
	t.Cleanup(db.Close)
	return NewDispatcher(db)
}

func handle(t *testing.T, d *Dispatcher, args ...string) string {
	t.Helper()
	return string(d.Handle(args))
}

func parseIntegerReply(t *testing.T, reply string) int64 {
	t.Helper()
	trimmed := strings.TrimPrefix(strings.TrimSuffix(reply, "\r\n"), ":")
	if trimmed == reply {
		t.Fatalf("reply %q is not a RESP integer", reply)
	}
	n, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		t.Fatalf("reply %q: %v", reply, err)
	}
	return n
}
