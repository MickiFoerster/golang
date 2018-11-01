package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type date struct {
	Year  int
	Month int
	Day   int
}

type MotoCrossDriver struct {
	Forename     string
	Surname      string
	Birthdate    date `json:"date_of_birth"`
	HeightInCm   int  `json:"height_in_cm,omitempty"`
	Currentbrand string
}

var drivers = []MotoCrossDriver{
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
	fmt.Printf("%s\n", data)
	// Print in pretier with indentation
	data, err = json.MarshalIndent(drivers, "", "    ")
	if err != nil {
		log.Fatalf("error: JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}
