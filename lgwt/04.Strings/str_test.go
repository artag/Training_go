package stringsResearch

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

// func Compare(a, b string) int

func ExampleCompare1() {
	fmt.Println(strings.Compare("ABC", "abc"))

	// Output:
	// -1
}

func ExampleCompare2() {
	fmt.Println(strings.Compare("123", "123"))

	// Output:
	// 0
}

func ExampleCompare3() {
	fmt.Println(strings.Compare("pax", "max"))

	// Output:
	// 1
}

func ExampleCompare4() {
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

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func ExampleCount() {
	fmt.Println(strings.Count("Alabama", "a"))
	fmt.Println(strings.Count("cheese", "e"))
	fmt.Println(strings.Count("five", ""))  // before & after each rune (4 + 1)
	fmt.Println(strings.Count("seven", "")) // 6

	// Output:
	// 3
	// 3
	// 5
	// 6
}

// func Cut(s, sep string) (before, after string, found bool)

// Cut slices s around the first instance of sep, returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
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

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding,
// which is a more general form of case-insensitivity.
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

// Fields splits the string s around each instance of one or more
// consecutive white space characters, as defined by unicode.IsSpace,
// returning a slice of substrings of s
// or an empty slice if s contains only white space.
func ExampleFields() {
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))

	// Output:
	// Fields are: ["foo" "bar" "baz"]
}

// func FieldsFunc(s string, f func(rune) bool) []string

// FieldsFunc splits the string s at each run of Unicode code points c
// satisfying f(c) and returns an array of slices of s.
// If all code points in s satisfy f(c) or the string is empty,
// an empty slice is returned.
func ExampleFieldsFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3 ban4...", f))

	// Output:
	// Fields are: ["foo1" "bar2" "baz3" "ban4"]
}

// func HasPrefix(s, prefix string) bool

// HasPrefix tests whether the string s begins with prefix.
func ExampleHasPrefix() {
	fmt.Println(strings.HasPrefix("Gopher", "Go"))
	fmt.Println(strings.HasPrefix("Gopher", "C"))
	fmt.Println(strings.HasPrefix("Gopher", ""))

	// Output:
	// true
	// false
	// true
}

// func HasSuffix(s, suffix string) bool

// HasSuffix tests whether the string s ends with suffix.
func ExampleHasSuffix() {
	fmt.Println(strings.HasSuffix("Amigo", "go"))
	fmt.Println(strings.HasSuffix("Amigo", "O"))
	fmt.Println(strings.HasSuffix("Amigo", "Ami"))
	fmt.Println(strings.HasSuffix("Amigo", ""))

	// Output:
	// true
	// false
	// false
	// true
}

// func Index(s, substr string) int

// Index returns the index of the first instance of substr in s,
// or -1 if substr is not present in s.
func ExampleIndex() {
	fmt.Println(strings.Index("chicken", "ken"))
	fmt.Println(strings.Index("chicken", "ch"))
	fmt.Println(strings.Index("chicken", "dmr"))

	// Output:
	// 4
	// 0
	// -1
}

// func IndexAny(s, chars string) int

// IndexAny returns the index of the first instance of any Unicode code point from chars in s,
// or -1 if no Unicode code point from chars is present in s.
func ExampleIndexAny() {
	fmt.Println(strings.IndexAny("chicken", "aeiouy"))
	fmt.Println(strings.IndexAny("crwth", "aeiouy"))

	// Output:
	// 2
	// -1
}

// func IndexByte(s string, c byte) int

// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
func ExampleIndexByte() {
	fmt.Println(strings.IndexByte("golang", 'g'))
	fmt.Println(strings.IndexByte("gophers", 'h'))
	fmt.Println(strings.IndexByte("golang", 'x'))

	// Output:
	// 0
	// 3
	// -1
}

// func IndexFunc(s string, f func(rune) bool) int

// IndexFunc returns the index into s of the first Unicode code point satisfying f(c),
// or -1 if none do.
func ExampleIndexFunc() {
	f := func(c rune) bool {
		return unicode.Is(unicode.Cyrillic, c)
	}
	fmt.Println(strings.IndexFunc("Hello, мир", f))
	fmt.Println(strings.IndexFunc("Hello, world", f))

	// Output:
	// 7
	// -1
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
