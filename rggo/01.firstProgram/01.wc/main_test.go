package main

import (
	"bytes"
	"testing"
)

// Tests the count function set to count words
func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4")
	exp := 4

	res := count(b)

	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}
