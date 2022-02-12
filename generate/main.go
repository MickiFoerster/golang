package main

import (
	"fmt"
)

//go:generate stringer -type=DirtbikeBrand
type DirtbikeBrand int

const (
	Honda DirtbikeBrand = iota
	Kawasaki
	Yamaha
	KTM
	GasGas
	Suzuki
	Husqwarna
)

func main() {
	var kawa DirtbikeBrand = Kawasaki
	fmt.Printf("Variable is of Brand %q\n", kawa)
}
