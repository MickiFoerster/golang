package main

import (
	"fmt"
	"log"

	redis "github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("error while connecting to redis server: %v", err)
	}
	defer conn.Close()

	reply, err := conn.Do("SET", "hello", "world!")
	if err != nil {
		log.Fatalf("error with command SET: %v", err)
	}

	reply, err = conn.Do("GET", "hello")
	if err != nil {
		log.Fatalf("error with command GET: %v", err)
	}

	switch v := reply.(type) {
	case nil:
		fmt.Printf("no value for given key\n")
	case []uint8:
		fmt.Printf("answer: %v\n", string(v))
	case string:
		fmt.Printf("answer: %v\n", v)
	default:
		fmt.Printf("unknown type: %T\n", v)
	}

}
