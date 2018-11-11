package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
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

const templ = `<html>
<body>
<ul>
{{range .}}
<li>{{.Forename}} {{.Surname}}</li>
  <ul>
		<li>Birthdate: {{.Birthdate.Day}}.{{.Birthdate.Month}}.{{.Birthdate.Year}}</li>
		<li>Height: {{.HeightInCm}} cm</li>
		<li>Brand: {{.Currentbrand}}</li>
  </ul>
{{end}}
</ul>
</body>
</html>
`

var listOfDrivers = template.Must(template.New("driversList").Parse(templ))

func main() {
	jsonFile, err := os.Open("drivers.json")
	if err != nil {
		log.Fatalln("error: Cannot open JSON file: ", err)
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("error: Could not read JSON data: ", err)
	}
	var drivers []MotoCrossDriver
	if err = json.Unmarshal(jsonData, &drivers); err != nil {
		log.Fatalln("error: Could not unmarshal JSON data: ", err)
	}
	fmt.Println(drivers)

	if err := listOfDrivers.Execute(os.Stdout, drivers); err != nil {
		log.Fatal("error: Could not execute template: ", err)
	}
}
