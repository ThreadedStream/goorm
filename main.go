package main

import (
	"log"
)

func main() {
	contents, err := readRaw("schema.sch")
	if err != nil {
		log.Fatal(err)
	}
	err = parseSchema(contents)
	if err != nil {
		log.Fatal(err)
	}

	conn := Conn{
		connString: "host=127.0.0.1 port=5432 user=postgres password=135797531 dbname=postgres sslmode=disable",
	}

	err = conn.initialize()
	if err != nil {
		log.Fatal(err)
	}

}
