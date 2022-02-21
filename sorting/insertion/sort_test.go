package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

var SizeTestData = 1024 * 256
var testData = []uint8{}

func BenchmarkInsertion_search_sort_with_binary_search(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createTestdata(SizeTestData)

		v := []uint8{}
		v = append(v, testData...)

		insertion_search_sort_with_binary_search(v)
	}
}

func BenchmarkInsertion_sort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		createTestdata(SizeTestData)

		v := []uint8{}
		v = append(v, testData...)

		insertion_sort(v)
	}
}

func TestInsertion_search_sort_with_binary_search(t *testing.T) {
	fmt.Println("Creating test data ...")
	createTestdata(SizeTestData)

	v := []uint8{}
	v = append(v, testData...)

	fmt.Println("Now sort data ...")
	insertion_search_sort_with_binary_search(v)

	fmt.Println("Testing if sort criteria holds ...")
	for i := 0; i+1 < len(v); i++ {
		if v[i] > v[i+1] {
			t.Fatalf("Error at index %v: %v > %v. Not sorted.", i, v[i], v[i+1])
		}
	}
}

func TestInsertion_sort(t *testing.T) {
	fmt.Println("Creating test data ...")
	createTestdata(SizeTestData)

	v := []uint8{}
	v = append(v, testData...)

	fmt.Println("Now sort data ...")
	insertion_sort(v)

	fmt.Println("Testing if sort criteria holds ...")
	for i := 0; i+1 < len(v); i++ {
		if v[i] > v[i+1] {
			t.Fatalf("Error at index %v: %v > %v. Not sorted.", i, v[i], v[i+1])
		}
	}
}

func createTestdata(testSize int) {
	if len(testData) == testSize {
		return
	}

	fmt.Println("create new test data")
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v := []uint8{}

	reader := bufio.NewReader(f)
	for {
		buf := make([]uint8, 4096)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		buf = buf[:n]
		v = append(v, buf...)
		if len(v) > testSize {
			v = v[:testSize]
			break
		} else if len(v) == testSize {
			break
		}
	}

	testData = v
}
