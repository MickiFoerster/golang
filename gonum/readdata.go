package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	xys, err := readData("data.txt")
	if err != nil {
		log.Fatalf("Could not read data: %v", err)
	}
	fmt.Println(xys)
}

type xy struct{ x, y float64 }

func readData(path string) ([]xy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var xys []xy
	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("Discarding bad data point %q: %v", s.Text(), err)
			continue
		}
		xys = append(xys, xy{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("Could not read data: %v", err)
	}

	return xys, nil
}
