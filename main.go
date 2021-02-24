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

	var fields = make([]string, 0)
	fields = append(fields, "id", "name", "surname")

	var values = make([]interface{}, 0)
	values = append(values, 24, "Arno", "Dorian")

	err = conn.insertValues("person", fields, values)
	if err != nil {
		log.Fatal(err)
	}

	//val := reflect.ValueOf(&person).Elem()
	//x  := val.Field(0).String()
	//fmt.Printf("Value of the first field is %v, associated struct tag is %s\n", x, val.Type().Field(0).Tag.Get("db"))
}
