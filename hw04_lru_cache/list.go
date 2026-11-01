package hw04lrucache

type List[T any] interface {
	Len() int
	Front() *ListItem[T]
	Back() *ListItem[T]
	PushFront(v T) *ListItem[T]
	PushBack(v T) *ListItem[T]
	Remove(i *ListItem[T])
	MoveToFront(i *ListItem[T])
}

type ListItem[T any] struct {
	Value T
	Next  *ListItem[T]
	Prev  *ListItem[T]
}

type list[T any] struct {
	head   *ListItem[T]
	tail   *ListItem[T]
	length int
}

func (l *list[T]) Len() int {
	return l.length
}

func (l *list[T]) Front() *ListItem[T] {
	return l.head
}

func (l *list[T]) Back() *ListItem[T] {
	return l.tail
}

func (l *list[T]) PushFront(v T) *ListItem[T] {
	item := ListItem[T]{Value: v}
	if l.length == 0 {
		l.head, l.tail = &item, &item
	} else {
		l.head.Prev = &item
		item.Next = l.head
		l.head = &item
	}
	l.length++
	return &item
}

func (l *list[T]) PushBack(v T) *ListItem[T] {
	item := ListItem[T]{Value: v}
	if l.length == 0 {
		l.head, l.tail = &item, &item
	} else {
		l.tail.Next = &item
		item.Prev = l.tail
		l.tail = &item
	}
	l.length++
	return &item
}

func (l *list[T]) Remove(i *ListItem[T]) {
	if i != nil {
		switch {
		case i.Prev == nil && i.Next == nil:
			l.head, l.tail = nil, nil
		case i.Prev == nil:
			l.head = i.Next
			i.Next.Prev = nil
			i.Next = nil
		case i.Next == nil:
			l.tail = i.Prev
			i.Prev.Next = nil
			i.Prev = nil
		default:
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
			i.Prev, i.Next = nil, nil
		}
		l.length--
	}
}

func (l *list[T]) MoveToFront(i *ListItem[T]) {
	if i != nil {
		switch {
		case i.Prev == nil:
		case i.Next == nil:
			i.Prev.Next = nil
			l.tail = i.Prev
			i.Next = l.head
			l.head.Prev = i
			l.head = i
			i.Prev = nil
		default:
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
			i.Next = l.head
			l.head.Prev = i
			l.head = i
			i.Prev = nil
		}
	}
}

func (l *list[T]) TransformToSlice() []T {
	s := make([]T, 0, l.length)
	i := l.head
	for range l.length {
		s = append(s, i.Value)
		i = i.Next
	}
	return s
}

func NewList[T any]() List[T] {
	return new(list[T])
}
