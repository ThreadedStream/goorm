package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type SchemaToken int

const (
	SIDENTIFIER SchemaToken = 0x0
	STRING                  = 0x1
	INT                     = 0x2
	FLOAT                   = 0x3
	NOTNULL                 = 0x4
	LENGTH                  = 0x5
	PK                      = 0x6
	TABLE                   = 0x7
	SERIAL                  = 0x8
)

var lineNum = 1

const (
	COLUMN_NAME  string = "column_name"
	TABLE_NAME   string = "table_name"
	COLUMN_PROPS string = "column_props"
)

/*
	Parsed version:
		[table_name] = {
					column_name1: [properties of column_name1...]
	                column_name2: [properties of column_name2...]
                  }
*/

var statements = make(map[string]map[string][]interface{}, 0)

var tokensToStrings = map[SchemaToken]string{
	SIDENTIFIER: "^[_a-z]\\w*$",
	STRING:      "string=[0-9]+",
	INT:         "int",
	FLOAT:       "float",
	NOTNULL:     "!null",
	LENGTH:      "length=[0-9]+",
	PK:          "pk",
	TABLE:       "table=[a-zA-Z[_]*[0-9]*",
	SERIAL:      "serial",
}

func readRaw(path string) (string, error) {
	var bs []byte
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	bs, err = ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	contents := string(bs)
	return contents, nil
}

func parseSchema(contents string) error {
	blocks := strings.Split(contents, "\n\n")
	for _, block := range blocks {
		err := parseBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseBlock(contents string) error {
	currIndex := 0
	startIndex := 0

	//Parse table name
	if contents[0] == 0xA {
		currIndex++
	}
	for contents[currIndex] != 0xA {
		currIndex++
	}
	err, tableName := parseTableName(contents[0:currIndex])
	if err != nil {
		return err
	}
	for currIndex < len(contents) {
		startIndex = currIndex
		currIndex++
		for contents[currIndex] != 0xA {
			currIndex++
			if currIndex == len(contents) {
				break
			}
		}

		err = parseLine(contents[startIndex:currIndex], tableName)
		if err != nil {
			return err
		}
		startIndex = currIndex
		currIndex++
		lineNum++
		continue
	}
	return nil
}

func parseTableName(contents string) (error, string) {
	split := strings.Split(contents, "=")
	if strings.TrimSpace(split[0]) != "table" {
		return errors.New(fmt.Sprintf("Expected 'table', found %s", split[0])), ""
	}
	split[1] = strings.TrimSpace(split[1])
	statements[split[1]] = make(map[string][]interface{}, 0)
	return nil, split[1]
}

func parseLine(line string, tableName string) error {
	currIndex := 0
	currToken := ""
	parenStart := -1
	parenEnd := -1

	for currIndex != len(line)-1 {
		for line[currIndex] != ' ' && line[currIndex] != '(' {
			currToken += string(line[currIndex])
			currIndex++
			if currIndex == len(line)-1 {
				return errors.New(fmt.Sprintf("Expected '(' or ' ' at line %d\n", lineNum))
			}
		}
		columnName := strings.TrimSpace(currToken)
		currToken = ""
		if line[currIndex] == ' ' && line[currIndex+1] == '(' {
			currIndex++
			parenStart = currIndex
			currIndex++
		} else if line[currIndex] == '(' {
			parenStart = currIndex
			currIndex++
		} else {
			return errors.New(fmt.Sprintf("Expected '(' at line %d\n", lineNum))
		}
		for line[currIndex] != ')' {
			if currIndex == len(line)-1 && line[currIndex] != ')' {
				return errors.New(fmt.Sprintf("Expected '(' at line %d\n", lineNum))
			}
			currIndex++
		}
		parenEnd = currIndex
		err := parseParentheses(line[parenStart+1:parenEnd], tableName, columnName)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseParentheses(contents string, tableName string, columnName string) error {
	params := strings.Split(contents, ",")
	prettifyParams(params)
	for _, token := range params {
		switch token {
		case "string":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[STRING])
			break
		case "int":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[INT])
			break
		case "float":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[FLOAT])
			break
		case "!null":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[NOTNULL])
			break
		case "pk":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[PK])
			break
		case "serial":
			statements[tableName][columnName] = append(statements[tableName][columnName], tokensToStrings[SERIAL])
		default:
			matchesString, _ := regexp.MatchString(tokensToStrings[STRING], token)
			if matchesString {
				stringToken := make([]string, 0)
				strs := strings.Split(token, "=")
				stringToken = append(stringToken, strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]))
				statements[tableName][columnName] = append(statements[tableName][columnName], stringToken)
			} else {
				return errors.New(fmt.Sprintf("Unexpected identifier %s at line %d\n", token, lineNum))
			}
		}
	}

	return nil
}

func prettifyParams(params []string) {
	for i := 0; i < len(params); i++ {
		params[i] = strings.TrimSpace(params[i])
	}
}
