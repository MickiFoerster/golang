package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connection via Unix socket with path postgres:///dbname?host=/var/run/postgresql/
	db, err := sql.Open("postgres", "postgres://user1:password@localhost/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to Postgres database successful.")
}
