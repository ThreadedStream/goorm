package main

import (
	"fmt"
	"log"
)


func main() {
	contents, err := readRaw("schema.sch")
	if err != nil {
		log.Fatal(err)
	}
	parseSchema(contents)
	for line, contents := range statements {
		fmt.Printf("%d - %s\n", line, contents)
	}
}
