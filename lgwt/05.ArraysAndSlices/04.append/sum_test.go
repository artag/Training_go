package sum_slice

import (
	"reflect"
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

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}
	assertSlices(t, got, want)
}

func TestSumAllTails(t *testing.T) {
	t.Run(
		"make the sums of some slices 1",
		func(t *testing.T) {
			got := SumAllTails([]int{1, 2}, []int{0, 9})
			want := []int{2, 9}
			assertSlices(t, got, want)
		})
	t.Run(
		"make the sums of some slices 2",
		func(t *testing.T) {
			got := SumAllTails([]int{1, 2, 3}, []int{1, 9, 7})
			want := []int{5, 16}
			assertSlices(t, got, want)
		})
	t.Run(
		"safely sum empty slices",
		func(t *testing.T) {
			got := SumAllTails([]int{}, []int{3, 4, 5})
			want := []int{0, 9}
			assertSlices(t, got, want)
		})
}

func assertIntegers(t *testing.T, got, want int, numbers []int) {
	if got == want {
		return
	}
	t.Errorf("got %d want %d given, %v", got, want, numbers)
}

func assertSlices(t *testing.T, got, want []int) {
	if reflect.DeepEqual(got, want) {
		return
	}

	t.Errorf("got %v want %v", got, want)
}
