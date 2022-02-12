package main

import "fmt"

type Node struct {
	val  string
	prev *Node
	next *Node
}

type List struct {
	root *Node
}

func (l *List) addNode(val string) {
	if l.root == nil {
		l.root = &Node{val, nil, nil}
		return
	}

	// Find position where new element to insert
	t := l.root
	for t.val <= val && t.next != nil {
		t = t.next
	}

	if t.next == nil && t.val <= val {
		t.next = &Node{val, t, nil}
		return
	}

	newnode := &Node{val, t.prev, t}
	if t != l.root {
		t.prev.next = newnode
	} else {
		l.root = newnode
	}
	t.prev = newnode
}

func (l List) dump() {
	t := l.root
	for t != nil {
		fmt.Printf("%p -> %v\n", t, *t)
		t = t.next
	}
}

func main() {
	var l List
	l.addNode("b")
	l.addNode("a")

	l.dump()
}
