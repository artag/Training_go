package shapes

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {
	expected := 30.0
	actual := Perimeter(10.0, 5.0)
	assertFloat64(t, expected, actual)
}

func TestArea(t *testing.T) {
	expected := 72.0
	actual := Area(12.0, 6.0)
	assertFloat64(t, expected, actual)
}

func assertFloat64(t *testing.T, want float64, got float64) {
	if math.Abs(got-want) < math.SmallestNonzeroFloat64 {
		return
	}
	t.Errorf("want: %.2f got %.2f ", want, got)
}
