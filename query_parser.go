package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

type Conn struct {
	DB         *sql.DB //Database connection object
	connString string  //Connection string
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

	tableQuery += ");"
	return tableQuery, nil
}
