package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

//go:embed file.txt
var s string

//go:embed users.json
var users_data_json []byte

//go:embed users.yaml
var users_data_yaml []byte

type User struct {
	User  string `yaml:user,json:user`
	Email string `yaml:email,json:email`
}
type Userlist struct {
	Users []User `yaml:users,json:users`
}

func main() {
	// Print content of file.txt
	fmt.Print(s)

	// Show JSON file
	showJSON()

	// Show Yaml file
	showYAML()
}

func showYAML() {
	var users Userlist
	if err := yaml.Unmarshal(users_data_yaml, &users); err != nil {
		log.Fatal(err)
	}
	fmt.Println("YAML:", users)
}

func showJSON() {
	var users Userlist
	if err := json.Unmarshal(users_data_json, &users); err != nil {
		log.Fatal(err)
	}
	fmt.Println("JSON:", users)
}
