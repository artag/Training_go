package shapes

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {
	t.Run(
		"Find rectangle area",
		func(t *testing.T) {
			rectangle := Rectangle{10.0, 5.0}
			expected := 30.0
			actual := rectangle.Perimeter()
			assertFloat64(t, expected, actual)
		})
	t.Run(
		"Find circle area",
		func(t *testing.T) {
			circle := Circle{10.0}
			expected := 62.83185307179586
			actual := circle.Perimeter()
			assertFloat64(t, expected, actual)
		})
}

func TestArea(t *testing.T) {
	t.Run(
		"Find rectangle area",
		func(t *testing.T) {
			rectangle := Rectangle{12.0, 6.0}
			expected := 72.0
			actual := rectangle.Area()
			assertFloat64(t, expected, actual)
		})
	t.Run(
		"Find circle area",
		func(t *testing.T) {
			circle := Circle{10}
			expected := 314.1592653589793
			actual := circle.Area()
			assertFloat64(t, expected, actual)
		})

}

func assertFloat64(t *testing.T, want float64, got float64) {
	if math.Abs(got-want) < math.SmallestNonzeroFloat64 {
		return
	}
	t.Errorf("want: %g got: %g", want, got)
}
