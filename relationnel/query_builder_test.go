package main

import (
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestQueryMatching(t *testing.T){
	var queryPattern = "INSERT INTO [a-zA-Z_0-9]+ [(]{1} ([a-zA-Z_0-9],)* [)]{1}"
	var query = "INSERT INTO table_1 ()"
	
	re := regexp.MustCompile(queryPattern)
	assert.Equal(t, re.MatchString(query), true)
}