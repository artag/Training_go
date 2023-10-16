package main

import (
	"bytes"
	"testing"
)

// Tests the count function set to count words
func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4")
	exp := 4

	act := count(b, countWords)

	Assert(t, exp, act)
}

// Tests the count function set to count lines
func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\nline2\nline3 word1")
	exp := 3

	act := count(b, countLines)

	Assert(t, exp, act)
}

// Tests the count function set to count bytes
func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4")
	exp := 23

	act := count(b, countBytes)

	Assert(t, exp, act)
}

func Assert(t *testing.T, exp int, act int) {
	if exp == act {
		return
	}

	t.Errorf("Expected %d, got %d instead.\n", exp, act)
}
