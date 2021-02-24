package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadRaw(t *testing.T) {
	path := "schema.sch"

	_, err := readRaw(path)

	assert.Equal(t, err, nil)
}

func TestParseSchema(t *testing.T) {
	path := "schema.sch"
	contents, err := readRaw(path)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, contents, "")

	err = parseSchema(contents)

	assert.Equal(t, err, nil)
}
