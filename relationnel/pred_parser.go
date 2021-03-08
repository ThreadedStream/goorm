package main

/*
	Predicate language for building complex queries.
	Here is a list of tokens with their respective sql alternatives:

	General structure:
	Predicate token -> SQL token
	1) FILTER -> SELECT
	2) * -> *
	3) | -> WHERE
	4) == -> =
	5) ORDERING (column names...) -> ORDER BY
	6) + -> ASC
	7) - -> DESC

	Coming soon...
*/

type PredicateToken int

const (
	PIDENTIFIER PredicateToken = 0x0
	ALL                        = 0x1
	WHERE                      = 0x2
	EQUAL                      = 0x3
	GREATERTHAN                = 0x4
	LESSTHAN                   = 0x5
	ORDERBY                    = 0x6
	ASCENDING                  = 0x7
	DESCENDING                 = 0x8
)

var predicateTokensToStrings = map[PredicateToken]string{
	PIDENTIFIER: "^[_a-z]\\w*$",
	ALL:         "*",
	WHERE:       "|",
	EQUAL:       "==",
	GREATERTHAN: ">",
	LESSTHAN:    "<",
	ORDERBY:     "ORDERING ([^,]+)",
	ASCENDING:   "+",
	DESCENDING:  "-",
}

var predicateToSql = map[PredicateToken]string{
	PIDENTIFIER: "^[_a-z]\\w*$",
	ALL:         "*",
	WHERE:       "|",
	EQUAL:       "==",
	GREATERTHAN: ">",
	LESSTHAN:    "<",
	ORDERBY:     "ORDERING ([^,]+)",
	ASCENDING:   "+",
	DESCENDING:  "-",
}

func splitPredicateIntoTokens(predicate string) {
	var currIndex = 0

	var currToken []rune
	for currIndex < len(predicate) {
		for predicate[currIndex] != ' ' {
			currToken = append(currToken, rune(predicate[currIndex]))
		}
		switch string(currToken) {
		case "*":

		}
	}
}
