package store

import "github.com/moges7624/gredis/internal/datastruct/list"

type ListValue struct {
	List list.RedisList
}

func (l *ListValue) Type() ValueType {
	return ListType
}

func NewRedisList() list.RedisList {
	return list.NewLinkedList()
}

func (s *Store) LLen(key string) (int, error) {
	s.mu.RLock()
	e, exists := s.get(key)
	s.mu.RUnlock()

	if !exists {
		return 0, nil
	}

	lv, ok := e.value.(*ListValue)
	if !ok {
		return 0, ErrWrongType
	}

	return lv.List.Len(), nil
}

func (s *Store) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.get(key)

	if !exists {
		l := NewRedisList()
		l.RPush(values...)

		s.data[key] = entry{
			value: &ListValue{
				List: l,
			},
		}

		return len(values), nil
	}

	lv, ok := e.value.(*ListValue)
	if !ok {
		return 0, ErrWrongType
	}

	lv.List.RPush(values...)

	return lv.List.Len(), nil
}
