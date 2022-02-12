package main

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestList(t *testing.T) {
	test_data := []byte{'b', 'c', 'd', 'e', 'g', 'f', 'a', 'h'}
	var l List
	for _, d := range test_data {
		fmt.Print("Insert ", string(d), " into list ...")
		l.addNode(string(d))
		fmt.Println("done")
	}

	l.dump()

	// Checking sorted order
	elem := l.root
	last := elem
	for i := 0; i < len(test_data); i++ {
		// check value
		expected := string(byte(i + 'a'))
		assert.Equal(t, elem.val, expected, "Value should be equal")

		// check pointers
		switch i {
		case 0:
			if elem.prev != nil {
				t.Errorf("Previous pointer of first element must be nil")
			}
		case len(test_data) - 1:
			if elem.next != nil {
				t.Errorf("Pointer 'next' of last element must be nil")
			}
			assert.Equal(t, last, elem.prev, "Predecessor of last element points to second last element")
		default:
			assert.Equal(t, last.next, elem, "Successor of last element points to current element")
			assert.Equal(t, last, elem.prev, "Predecessor of current element points to previous element")
		}

		last = elem
		elem = elem.next
	}

}
