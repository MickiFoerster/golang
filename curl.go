package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func pretty_print_json(data []byte) {
	// response should be JSON format
	var json_data map[string]interface{}

	err := json.Unmarshal(data, &json_data)
	if err != nil {
		log.Fatalf("parsing body as JSON failed: %s", err)
	}

	// Access JSON data which is now mapped
	fmt.Println(json_data["origin"])
	if sub_obj, ok := json_data["headers"].(map[string]interface{}); ok {
		fmt.Println(sub_obj["User-Agent"])
	} else {
		log.Fatalln("could not read JSON sub object")
	}

	// Pretty printing of response
	formatted_data, err := json.MarshalIndent(json_data, "", " ")
	if err != nil {
		log.Fatalf("failed to transform JSON into formatted string: %s", err)
	}
	fmt.Println(string(formatted_data))
}

func main() {
	client := &http.Client{}

	fmt.Printf("Now we do a GET request ...\n")
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatalf("creating GET request failed: %s", err)
	}

	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("GET request failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("reading response body failed: %s", err)
	}

	pretty_print_json(body)

	fmt.Println("========================================================\n\n")
	type Response struct {
		Origin  string `json:"origin"`
		Url     string `json:"url"`
		Headers struct {
			Useragent string `json:"User-Agent"`
		} `json:"headers"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	fmt.Println("========================================================\n\n")

	fmt.Printf("Now we do a POST request ...\n")
	var post_data = strings.NewReader(`{"amount": "10.00", "description": "test"}`)
	req, err = http.NewRequest("POST", "https://httpbin.org/post", post_data)
	if err != nil {
		log.Fatalf("creating GET request failed: %s", err)
	}

	req.Header.Set("accept", "application/json")
	req.SetBasicAuth("USER", "password")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("GET request failed: %s", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("reading response body failed: %s", err)
	}

	pretty_print_json(body)

}
