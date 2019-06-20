package models

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}
