package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
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

/*
	"tableName" designates a name of the table which is to be populated with new data
	"fields" designates a list of column names, whereas "values" is respective list of values
	Beware that order matters.
*/
func (c *Conn) insertValues(tableName string, fields []string, values []interface{}) error {

	if len(values)%len(fields) != 0 {
		return errors.New("Length of fields must be the multiple of length of values\n")
	}
	var query = fmt.Sprintf("INSERT INTO %s (", tableName)

	for _, field := range fields {
		query += field + ","
	}

	//Omit unnecessary comma at the end
	query = string([]rune(query)[0:len(query)-1]) + ")"

	query += " VALUES ("
	for _, value := range values {
		stringValue, ok := value.(string)
		if ok {
			query += "'" + stringValue + "',"
		} else {
			query += fmt.Sprintf("%v,", value)
		}
	}

	//Omit unnecessary comma at the end
	query = string([]rune(query)[0:len(query)-1]) + ");"

	//Query execution stage
	_, err := c.DB.Query(query)
	if err != nil {
		return err
	}

	log.Println("Values were successfully placed into table")

	return nil
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

func (c *Conn) commit() (map[string]string, error) {
	var query = ""
	if len(c.queryPrims.selectedFields) > 0 {
		for _, field := range c.queryPrims.selectedFields {
			query += fmt.Sprintf("SELECT (%s,", field)
		}
		query = string([]rune(query[0:len(query)-1])) + ") "
	}

	if c.queryPrims.tableName != "" {
		query += "FROM " + c.queryPrims.tableName
	}

	if c.queryPrims.wherePred != "" {
		query += "WHERE " + c.queryPrims.wherePred
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

	var rowsPretty = prettifyRows(rows)

	return rowsPretty, nil
}

func prettifyRows(rows *sql.Rows) map[string]string {
	return nil
}
