package sum_slice

import (
	"testing"
)

func TestSum(t *testing.T) {
	t.Run(
		"collection of any size",
		func(t *testing.T) {
			numbers := []int{1, 2, 3, 4, 5}
			got := Sum(numbers)
			want := 15
			assertIntegers(t, got, want, numbers[:])
		})
}

func assertIntegers(t *testing.T, got, want int, numbers []int) {
	if got == want {
		return
	}
	t.Errorf("got %d want %d given, %v", got, want, numbers)
}
