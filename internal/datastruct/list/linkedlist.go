package list

import "container/list"

type LinkedList struct {
	list list.List
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		list: *list.New(),
	}
}

func (l *LinkedList) Len() int {
	return l.list.Len()
}

func (l *LinkedList) RPush(values ...string) int {
	for _, v := range values {
		l.list.PushFront(v)
	}

	return len(values)
}
