package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type date struct {
	Year  int
	Month int
	Day   int
}

type motoCrossDriver struct {
	Forename     string
	Surname      string
	Birthdate    date `json:"date_of_birth"`
	HeightInCm   int  `json:"height_in_cm,omitempty"`
	Currentbrand string
}

var drivers = []motoCrossDriver{
	{
		Forename: "Ken",
		Surname:  "Roczen",
		Birthdate: date{
			Year:  1994,
			Month: 4,
			Day:   29,
		},
		Currentbrand: "Honda Racing Corporation",
		HeightInCm:   170,
	},
	{
		Forename: "Eli",
		Surname:  "Tomac",
		Birthdate: date{
			Year:  1992,
			Month: 11,
			Day:   14,
		},
		Currentbrand: "Monster Energy Kawasaki",
		HeightInCm:   175,
	},
}

func main() {
	data, err := json.Marshal(drivers)
	if err != nil {
		log.Fatalf("error: JSON marshaling failed: %s", err)
	}
	writeFile("drivers.json", data)

	// Print in prettier with indentation
	data, err = json.MarshalIndent(drivers, "", "    ")
	if err != nil {
		log.Fatalf("error: JSON marshaling failed: %s", err)
	}
	writeFile("drivers-prettier.json", data)
}

func writeFile(fn string, data []byte) {
	fp, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	n, err := fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File %q (%d bytes) written successfully.\n", fn, n)
}
