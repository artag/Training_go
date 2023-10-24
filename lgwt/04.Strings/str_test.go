package stringsResearch

import (
	"fmt"
	"strings"
	"testing"
)

// func Compare(a, b string) int

func TestCompare(t *testing.T) {
	t.Run(
		"Compare a < b, result -1",
		func(t *testing.T) {
			a := "ABC"
			b := "abc"
			expected := -1

			actual := strings.Compare(a, b)

			AssertInt(t, expected, actual)
		})
	t.Run(
		"Compare a == b, result 0",
		func(t *testing.T) {
			a := "123"
			b := "123"
			expected := 0

			actual := strings.Compare(a, b)

			AssertInt(t, expected, actual)
		})
	t.Run(
		"Compare a > b, result 1",
		func(t *testing.T) {
			a := "pax"
			b := "max"
			expected := 1

			actual := strings.Compare(a, b)

			AssertInt(t, expected, actual)
		})
}

func ExampleCompare() {
	fmt.Println(strings.Compare("a", "b"))
	fmt.Println(strings.Compare("a", "a"))
	fmt.Println(strings.Compare("b", "a"))

	// Output:
	// -1
	// 0
	// 1
}

// func Contains(s, substr string) bool

func TestContains(t *testing.T) {
	t.Run(
		"Contains new line symbol",
		func(t *testing.T) {
			str := "some string\n"
			actual := strings.Contains(str, "\n")
			AssertTrue(t, actual)
		})

	t.Run(
		"Contains operation is case sensitive",
		func(t *testing.T) {
			str := "SOME STRING"
			actual := strings.Contains(str, "rin")
			AssertFalse(t, actual)
		})
}

func ExampleContains() {
	fmt.Println(strings.Contains("seafood", "foo"))
	fmt.Println(strings.Contains("seafood", "bar"))
	fmt.Println(strings.Contains("seafood", ""))
	fmt.Println(strings.Contains("", ""))

	// Output:
	// true
	// false
	// true
	// true
}

// func ContainsAny(s, chars string) bool

func ExampleContainsAny() {
	fmt.Println(strings.ContainsAny("team", "i"))
	fmt.Println(strings.ContainsAny("fail", "ui"))
	fmt.Println(strings.ContainsAny("ure", "ui"))
	fmt.Println(strings.ContainsAny("failure", "ui"))
	fmt.Println(strings.ContainsAny("foo", ""))
	fmt.Println(strings.ContainsAny("", ""))

	// Output:
	// false
	// true
	// true
	// true
	// false
	// false
}

// Asserts

func AssertInt(t *testing.T, expected int, actual int) {
	if expected == actual {
		return
	}
	t.Errorf("\n"+
		"Expected: %d\n"+
		"Actual: %d\n",
		expected, actual)
}

func AssertTrue(t *testing.T, actual bool) {
	if actual {
		return
	}

	t.Errorf("Error: actual result is false\n")
}

func AssertFalse(t *testing.T, actual bool) {
	if !actual {
		return
	}

	t.Errorf("Error: actual result is true\n")
}
