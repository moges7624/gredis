package command

import (
	"strings"

	"github.com/moges7624/gredis/internal/resp"
	"github.com/moges7624/gredis/internal/store"
)

type handlerFunc func(s *store.Store, args []string) []byte

type Dispatcher struct {
	store    *store.Store
	handlers map[string]handlerFunc
}

func NewDispatcher(s *store.Store) *Dispatcher {
	d := &Dispatcher{
		store:    s,
		handlers: make(map[string]handlerFunc),
	}

	d.register()
	return d
}

func (d *Dispatcher) register() {
	d.handlers["PING"] = handlePing
	d.handlers["ECHO"] = handleEcho
	d.handlers["INFO"] = handleInfo

	d.handlers["TTL"] = handleTTL
	d.handlers["EXISTS"] = handleExists
	d.handlers["TYPE"] = handleType

	d.handlers["GET"] = handleGet
	d.handlers["SET"] = handleSet
	d.handlers["DEL"] = handleDel
	d.handlers["INCR"] = handleIncr
	d.handlers["DECR"] = handleDecr
}

func (d *Dispatcher) Handle(args []string) []byte {
	if len(args) == 0 {
		return resp.EncodeError("unknown command")
	}

	name := strings.ToUpper(args[0])
	h, ok := d.handlers[name]
	if !ok {
		return resp.EncodeError("unknown command '" + name + "'")
	}

	return h(d.store, args[1:])
}
