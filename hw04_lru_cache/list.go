package hw04lrucache

import (
    "fmt"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Print()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstNode *ListItem
	lastNode  *ListItem
	len       int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	if l.len > 0 {
		return l.firstNode
	}
	return nil
}

func (l list) Back() *ListItem {
	if l.len > 0 {
		return l.lastNode
	}

	return nil
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{v, nil, nil}

	l.len++

	if l.len > 1 {
		newItem.Next = l.Front()
		l.firstNode.Prev = &newItem
	}

	l.firstNode = &newItem

	if l.len == 1 {
		l.lastNode = &newItem
	}

	return l.Front()
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{v, nil, nil}

	l.len++

	if l.len > 1 {
		newItem.Prev = l.Back()
		l.lastNode.Next = &newItem
	}

	l.lastNode = &newItem

	if l.len == 1 {
		l.firstNode = &newItem
	}

	return l.Back()
}

func (l *list) Remove(i *ListItem) {
	l.len--

	if l.Front() == i {
		l.firstNode = i.Next
	}

	if l.Back() == i {
		l.lastNode = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Front() == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.lastNode = i.Prev
	}

	i.Prev = nil
	i.Next = l.Front()
	l.Front().Prev = i
	l.firstNode = i
}

func (l list) Print() {
	if l.firstNode == nil {
		fmt.Println(nil)
		return
	}

	i := l.Front()

	for i != nil {
		fmt.Println(i.Value)
		i = i.Next
	}
}