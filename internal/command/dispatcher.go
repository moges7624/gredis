package command

import (
	"fmt"
	"strings"

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
	d.handlers["INFO"] = handleInfo

	d.handlers["GET"] = handleGet
	d.handlers["SET"] = handleSet
}

func (d *Dispatcher) Handle(args []string) []byte {
	if len(args) == 0 {
		return []byte("-ERR unknown command\r\n")
	}

	name := strings.ToUpper(args[0])
	h, ok := d.handlers[name]
	if !ok {
		return fmt.Appendf(nil, "-ERR unknown command '%s'\r\n", name)
	}

	return h(d.store, args[1:])
}
