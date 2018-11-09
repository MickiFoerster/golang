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

func main() {
	var motoCrossDriver []MotoCrossDriver
	jsondata := []byte(`[{"Forename":"Ken","Surname":"Roczen","date_of_birth":{"Year":1994,"Month":4,"Day":29},"height_in_cm":170,"Currentbrand":"Honda Racing Corporation"},{"Forename":"Eli","Surname":"Tomac","date_of_birth":{"Year":1992,"Month":11,"Day":14},"height_in_cm":175,"Currentbrand":"Monster Energy Kawasaki"}]`)
	if err := json.Unmarshal(jsondata, &motoCrossDriver); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	for _, driver := range motoCrossDriver {
		fmt.Println(driver)
	}
}

func (mxd MotoCrossDriver) String() string {
	return fmt.Sprint(mxd.Forename, " ", mxd.Surname)
}
