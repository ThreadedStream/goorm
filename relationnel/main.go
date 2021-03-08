package main


import (
)

func add(a,b int) int{
	return a + b
}

func main() {
	
	var root = SampleTree()

	InsertKey(root, 4)

	PrintTree(root)

	// var root = SampleTree()
	
	// key := 4

	// InsertKey(root, key)
	//var node = SearchKey(root, key)


	// contents, err := readRaw("schema.sch")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = parseSchema(contents)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// conn := Conn{
	// 	connString: "host=127.0.0.1 port=5432 user=postgres password=135797531 dbname=postgres sslmode=disable",
	// }

	// err = conn.initialize()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var fields = make([]interface{}, 0)
	// fields = append(fields, "id", "name", "surname")

	// var values = make([]interface{}, 0)
	// values = append(values, 24, "Arno", "Dorian")

	// err = conn.insertValues("person", fields, values)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//var mapped = make([]map[string]interface{}, 0)
	//var tableName = "person"
	//var wherePred = "name = 'Arno' AND id > 1"
	//
	//mapped, err = conn.filter(fields).
	//	from(tableName).
	//	where(wherePred).
	//	commit()
	//
	//println(mapped)

	//val := reflect.ValueOf(&person).Elem()
	//x  := val.Field(0).String()
	//fmt.Printf("Value of the first field is %v, associated struct tag is %s\n", x, val.Type().Field(0).Tag.Get("db"))
}
