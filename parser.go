package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Tokens int

const (
	IDENTIFIER Tokens = 0x0
	STRING            = 0x1
	INT               = 0x2
	FLOAT             = 0x3
	NOTNULL           = 0x4
	LENGTH            = 0x5
	PK                = 0x6
	TABLE			  = 0x7
)

/*
	TODO: Transform input into the following format
	Parsed version: TABLE = <table_name>, COLUMN_NAME = <column_name>, COLUMN_PROPS = [properties of columns...]
 */

var statements = make(map[int][]string, 0)

var tokensToStrings = map[Tokens]string{
	IDENTIFIER: "^[_a-z]\\w*$",
	STRING:     "string",
	INT:        "int",
	FLOAT:      "float",
	NOTNULL:    "!null",
	LENGTH:     "length=[0-9]+",
	PK:         "pk",
	TABLE: 	    "table=[a-zA-Z[_]*[0-9]*",
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

func parseSchema(contents string) {
	currIndex := 0
	startIndex := 0

	currIndex++
	lineNum := 1
	for currIndex < len(contents) {
		for contents[currIndex] != 0xA {
			currIndex++
		}

		err := parseLine(contents[startIndex:currIndex], lineNum)
		if err != nil {
			log.Fatal(err)
		}
		startIndex = currIndex
		currIndex++
		lineNum++
		continue
	}
}

func parseLine(line string, lineNum int) error {
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
		statements[lineNum] = append(statements[lineNum], currToken)
		currToken = ""
		if line[currIndex] == ' ' && line[currIndex+1] == '('{
			currIndex++
			parenStart = currIndex
			currIndex++
		}else if line[currIndex] == '('{
			parenStart = currIndex
			currIndex++
		}else {
			return errors.New(fmt.Sprintf("Expected '(' at line %d\n", lineNum))
		}
		for line[currIndex] != ')' {
			if currIndex == len(line)-1 && line[currIndex] != ')' {
				return errors.New(fmt.Sprintf("Expected '(' at line %d\n", lineNum))
			}
			currIndex++
		}
		parenEnd = currIndex
		err := parseParentheses(line[parenStart+1:parenEnd], lineNum)
		if err != nil {
			return err
		}
	}

	return nil
}


func parseParentheses(contents string, lineNum int) error {
	params := strings.Split(contents, ",")
	prettifyParams(params)
	for _, token := range params {
		switch token {
		case "string":
			statements[lineNum] = append(statements[lineNum], tokensToStrings[STRING])
			break
		case "int":
			statements[lineNum] = append(statements[lineNum], tokensToStrings[INT])
			break
		case "float":
			statements[lineNum] = append(statements[lineNum], tokensToStrings[FLOAT])
			break
		case "!null":
			statements[lineNum] = append(statements[lineNum], tokensToStrings[NOTNULL])
			break
		case "pk":
			statements[lineNum] = append(statements[lineNum], tokensToStrings[PK])
			break
		default:
			matchesLength, _ := regexp.MatchString(tokensToStrings[LENGTH], token)
			matchesTable, _ := regexp.MatchString(tokensToStrings[TABLE], token)
			if matchesLength {
				statements[lineNum] = append(statements[lineNum], token)
			} else if matchesTable {
				statements[lineNum] = append(statements[lineNum], token)
			} else{
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
