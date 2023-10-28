package stringsResearch

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
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

// func ContainsRune(s string, r rune) bool

func ExampleContainsRune() {
	// Finds whether a string contains a particular Unicode code point.
	// The code point for the lowercase letter "a", for example, is 97.
	fmt.Println(strings.ContainsRune("aardvark", 97))
	fmt.Println(strings.ContainsRune("timeout", 97))

	// Output:
	// true
	// false
}

// func Count

func TestCount(t *testing.T) {
	t.Run(
		"count is case sensitive",
		func(t *testing.T) {
			actual := strings.Count("Alabama", "a")
			AssertInt(t, 3, actual)
		})
}

func ExampleCount() {
	fmt.Println(strings.Count("cheese", "e"))
	fmt.Println(strings.Count("five", ""))  // before & after each rune (4 + 1)
	fmt.Println(strings.Count("seven", "")) // 6

	// Output:
	// 3
	// 5
	// 6
}

// func Cut(s, sep string) (before, after string, found bool)

func ExampleCut() {
	show := func(s, sep string) {
		before, after, found := strings.Cut(s, sep)
		fmt.Printf("Cut(%q, %q) = %q, %q, %v\n", s, sep, before, after, found)
	}
	show("Gopher", "Go")
	show("Gopher", "ph")
	show("Gopher", "er")
	show("Gopher", "Badger")

	// Output:
	// Cut("Gopher", "Go") = "", "pher", true
	// Cut("Gopher", "ph") = "Go", "er", true
	// Cut("Gopher", "er") = "Goph", "", true
	// Cut("Gopher", "Badger") = "Gopher", "", false
}

// func EqualFold(s, t string) bool

func ExampleEqualFold() {
	fmt.Println(strings.EqualFold("ГО", "го"))
	fmt.Println(strings.EqualFold("go", "go"))
	fmt.Println(strings.EqualFold("Go", "go"))
	fmt.Println(strings.EqualFold("GO", "go"))
	fmt.Println(strings.EqualFold("go ", "go"))

	// Output:
	// true
	// true
	// true
	// true
	// false
}

// func Fields(s string) []string

func ExampleFields() {
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))

	// Output:
	// Fields are: ["foo" "bar" "baz"]
}

// func FieldsFunc(s string, f func(rune) bool) []string

func ExampleFieldsFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3 ban4...", f))

	// Output:
	// Fields are: ["foo1" "bar2" "baz3" "ban4"]
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
