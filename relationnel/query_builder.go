package main

import (
	"regexp"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"strconv"
	"strings"
)


type QueryPrimitives struct {
	selectedFields []string        // SELECT (selectedFields...)
	tableName      string          //FROM tableName
	wherePred      string          //WHERE wherePred
	orderBy        map[string]rune //ORDER BY [field:[+-]} | '+' - ascending order, '-' - descending order
}

type Conn struct {
	DB         *sql.DB         //Database connection object
	connString string          //Connection string
	queryPrims QueryPrimitives // Query primitives used for building filtering queries
}

//TODO: (Provide support for other drivers)
//NOTE: Supports only postgresql
func (c *Conn) initialize() error {
	var err error
	var created int
	if c.connString == "" {
		return errors.New("You must fill out connection string first\n")
	}

	c.DB, err = sql.Open("postgres", c.connString)
	if err != nil {
		return err
	}
	log.Println("Creating necessary tables...")
	created, err = c.createTables(statements)
	log.Printf("Successfully created %d tables\n", created)
	return nil
}

//Creates necessary tables defined in .sch file
func (c *Conn) createTables(statements map[string]map[string][]interface{}) (int, error) {
	created := 0
	for table, columns := range statements {
		query, err := createTableQuery(table, columns)
		if err != nil {
			return created, err
		}
		_, err = c.DB.Query(query)
		if err != nil {
			return created, err
		}
		created++
	}

	return created, nil
}

//Constructs a query for table creation
func createTableQuery(tableName string, columns map[string][]interface{}) (string, error) {
	tableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	for columnName, props := range columns {
		columnPropsString := columnName + " "
		//Parse type
		switch props[0] {
		case "int":
			columnPropsString += "INT "
			break
		case "float":
			columnPropsString += "FLOAT "
		default:
			list, ok := props[0].([]string)
			if ok {
				if list[0] == "string" {
					length, err := strconv.Atoi(list[1])
					if err != nil {
						return "", err
					}
					columnPropsString += fmt.Sprintf("VARCHAR ( %d ) ", length)
				}
			}
		}
		//Parse the rest
		for _, meta := range props[1 : len(props)-1] {
			switch meta {
			case "!null":
				columnPropsString += "NOT NULL "
				break
			case "pk":
				columnPropsString += "PRIMARY KEY "
			}
		}
		columnPropsString += ","
		tableQuery += columnPropsString
	}

	temp := []rune(tableQuery)
	tableQuery = string(temp[0 : len(temp)-2])
	tableQuery += ");"
	return tableQuery, nil
}

createQuery("; DROP TABLE table;")
func createQuery(table, field1, field2, field3 string){
	//INSERT INTO ;DROP TABLE table;
	var query = "INSERT INTO " + table + "(" + field1 + ","
}

/*
	"tableName" designates a name of the table which is to be populated with new data
	"fields" designates a list of column names, whereas "values" is respective list of values
	Beware that order matters.
*/
func (c *Conn) insertValues(tableName string, fields []interface{}, values []interface{}) error {

	if len(values)%len(fields) != 0 {
		return errors.New("Length of fields must be the multiple of length of values\n")
	}

	var questionMarks strings.Builder
	for i := 0; i < len(values); i++{
		questionMarks.WriteString("? ")
	}
	var queryPattern = "[INSERT|insert] [INTO|into] ^[a-zA-Z][_]*[0-9]+$ (^([a-zA-Z_0-9],)+$) VALUES (" + questionMarks.String() + ")"
	
	var query = "INSERT INTO table (id, name, surname) VALUES ($1, $2, $3)"

	var re = regexp.MustCompile(queryPattern)

	var query = "INSERT INTO " + tableName + " (" 

	for i := 0; i < len(fields); i++{
		strValue, ok := fields[i].(string); if ok{
			query += strValue
		}else{
			return errors.New("Field should be of type string\n")
		}
	}

	//Omit unnecessary comma at the end
	query = string([]rune(query)[0:len(query)-1]) + ") VALUES ("

	j := 0
	for i := 0; i < len(values) / len(values); i++{
		/*
			Account for the fact that there might be more than 1 insertion.
			For instance, consider the following scenario:
			INSERT INTO person(id, name, surname) VALUES (1, 'Arno', 'Dorian'), (2, 'Vito', 'Scaletta').
			In this case, i is running through all subsets of the whole set, whereas j is just an ordinary index
			of values list
		*/

		for j < len(values) + (i * len(values)){
			query += "?,"
			j++
		}
		query = string([]rune(query)[0:len(query)-1]) + ")"
	}


	if !re.MatchString(query){
		return errors.New("It seems like your input is invalid\n")
	}
	result, err := c.DB.Query(query); if err != nil{
		return err
	}


	log.Printf("Rows affected %v", result)

	return nil
}

func mergeLists(tableName string, src ... []interface{}) []interface{}{
	dest := make([]interface{}, 0)
	dest = append(dest, tableName)
	for i := 0; i < len(src); i++{
		for j := 0; j < len(src[i]); j++{
			dest = append(dest, src[i][j])
		}
	}

	return dest
}

//SELECT (id, name, surname) FROM person WHERE name = 'Arno'

func (c *Conn) filter(fields []string) *Conn {
	c.queryPrims.selectedFields = fields

	return c
}

func (c *Conn) from(tableName string) *Conn {
	c.queryPrims.tableName = tableName

	return c
}

func (c *Conn) where(wherePred string) *Conn {
	c.queryPrims.wherePred = wherePred

	return c
}

func (c *Conn) orderBy(fields map[string]rune) *Conn {
	c.queryPrims.orderBy = fields

	return c
}

func (c *Conn) commit() ([]map[string]interface{}, error) {
	var query = ""
	if len(c.queryPrims.selectedFields) > 0 {
		query += "SELECT ("
		for _, field := range c.queryPrims.selectedFields {
			query += fmt.Sprintf("%s,", field)
		}
		query = string([]rune(query[0:len(query)-1])) + ") "
	}

	if c.queryPrims.tableName != "" {
		query += "FROM " + c.queryPrims.tableName
	}

	if c.queryPrims.wherePred != "" {
		query += " WHERE " + c.queryPrims.wherePred
	}

	if len(c.queryPrims.orderBy) != 0 {
		query += "ORDER BY "
		for k, v := range c.queryPrims.orderBy {
			query += k
			if v == '+' {
				query += "ASC"
			} else if v == '-' {
				query += "DESC"
			}
		}
	}

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rowsPretty []map[string]interface{}

	//var person *Person

	//err = structify(rows, &person); if err != nil{
	//	return nil, err
	//}

	rowsPretty, err = prettifyRows(rows, c.queryPrims.selectedFields); if err != nil{
		return nil, err
	}

	return rowsPretty, nil
}

func structify(rows *sql.Rows, modelStruct interface{}) error{
	var modelPtr = reflect.ValueOf(modelStruct).Elem()
	var values = make([]interface{}, modelPtr.Type().Elem().NumField())

	for rows.Next(){
		//Allocate space for a pointer of modelElem's type
		var rowPtr = reflect.New(modelPtr.Type().Elem())
		//Get its value
		var rowValue = rowPtr.Elem()

		for i := 0; i < rowValue.NumField(); i++{
			values[i] = rowValue.Field(i).Addr().Interface()
		}
		err := rows.Scan(values...); if err != nil{
			return err
		}
		modelPtr.Set(reflect.Append(modelPtr, rowValue))
	}

	return nil
}

func prettifyRows(rows *sql.Rows, fields []string) ([]map[string]interface{},error) {
	var values []interface{}
	var mapped = make([]map[string]interface{}, 0)
	for rows.Next(){
		err := rows.Scan(pq.Array(&values)); if err != nil{
			return nil, err
		}

		var innerMap = make(map[string]interface{})
		for i, field := range fields{
			innerMap[field] = values[i]
		}
		mapped = append(mapped, innerMap)
	}

	return mapped, nil
}
