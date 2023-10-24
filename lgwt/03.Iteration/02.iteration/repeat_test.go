package iteration

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	actual := Repeat("a", 5)
	expected := "aaaaa"

	if expected != actual {
		t.Errorf("\n"+
			"Expected:\n"+
			"%q\n"+
			"Actual:\n"+
			"%q\n",
			expected, actual)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}

func ExampleRepeat() {
	example1 := Repeat("x", 3)
	example2 := Repeat("n", 7)

	fmt.Println(example1)
	fmt.Println(example2)

	// Output:
	// xxx
	// nnnnnnn
}
