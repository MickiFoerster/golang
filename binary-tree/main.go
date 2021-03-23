package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Tree struct {
	Left  *Tree
	ID    uint
	Value int
	Right *Tree
}

var counter uint = 0

func traverse(w io.Writer, t *Tree) {
	if t == nil {
		return
	}
	traverse(w, t.Left)
	fmt.Fprintf(w, `%v [label="%v"]`, t.ID, t.Value)
	fmt.Fprintln(w)
	if t.Left != nil {
		fmt.Fprintf(w, "%v -> %v\n", t.ID, t.Left.ID)
	}
	if t.Right != nil {
		fmt.Fprintf(w, "%v -> %v\n", t.ID, t.Right.ID)
	}
	traverse(w, t.Right)
}

func create(n int) *Tree {
	var t *Tree
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2*n; i++ {
		temp := rand.Intn(n * 2)
		t = insert(t, temp)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		counter += 1
		return &Tree{nil, counter, v, nil}
	}
	if v == t.Value {
		return t
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
		return t
	}
	t.Right = insert(t.Right, v)
	return t
}

func dump(w io.Writer, t *Tree) {
	fmt.Fprintln(w, "digraph {")
	traverse(w, t)
	fmt.Fprintln(w, "}")
}

func main() {
	tree := create(100)
	f, err := os.Create("bintree.dot")
	if err != nil {
		log.Fatalf("error: could not create file: %v\n", err)
	}
	dump(f, tree)
	f.Close()

	tree = insert(tree, -10)
	tree = insert(tree, 2)
	f, err = os.Create("bintree2.dot")
	if err != nil {
		log.Fatalf("error: could not create file: %v\n", err)
	}
	dump(f, tree)
	f.Close()
}
