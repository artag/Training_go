package iteration

import "testing"

func TestRepeat(t *testing.T) {
	actual := Repeat("a")
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
		Repeat("a")
	}
}
