package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	assert := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("\n"+
				"got: %q\n"+
				"want: %q\n",
				got, want)
		}
	}

	t.Run(
		"saying hello to people",
		func(t *testing.T) {
			got := Hello("Chris", "")
			want := "Hello, Chris"
			assert(t, got, want)
		})

	t.Run(
		"say 'Hello, World' when an empty string is supplied",
		func(t *testing.T) {
			got := Hello("", "")
			want := "Hello, World"
			assert(t, got, want)
		})

	t.Run(
		"in Spanish",
		func(t *testing.T) {
			got := Hello("Elodie", "Spanish")
			want := "Hola, Elodie"
			assert(t, got, want)
		})
}
