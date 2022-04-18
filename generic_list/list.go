package main

import (
	"fmt"
	"strings"
)

// List represents a singly-linked list that holds values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func (l *List[T]) String() string {
	var sb strings.Builder
	l.addVal(&sb)
	nxt := l.next
	for nxt != nil {
		sb.WriteString(", ")
		nxt.addVal(&sb)
		nxt = nxt.next
	}
	return sb.String()
}

func (l *List[T]) addVal(sb *strings.Builder) {
	str := fmt.Sprintf("%v", l.val)
	sb.WriteString(str)
}

func (prev *List[T]) Add(item T) *List[T] {
	next := List[T]{nil, item}
	prev.next = &next
	return &next
}

func main() {
	lst1 := List[int]{nil, 1}
	lst1.Add(2).Add(3)
	fmt.Println(lst1.String())

	lst2 := List[float64]{nil, 1.1}
	lst2.Add(1.2).Add(1.3)
	fmt.Println(lst2.String())
}
