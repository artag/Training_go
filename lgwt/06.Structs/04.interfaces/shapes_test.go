package shapes

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {
	t.Run(
		"Find rectangle perimeter",
		func(t *testing.T) {
			rectangle := Rectangle{10.0, 5.0}
			checkPerimeter(t, rectangle, 30.0)
		})
	t.Run(
		"Find circle perimeter",
		func(t *testing.T) {
			circle := Circle{10.0}
			checkPerimeter(t, circle, 62.83185307179586)
		})
}

func TestArea(t *testing.T) {
	t.Run(
		"Find rectangle area",
		func(t *testing.T) {
			rectangle := Rectangle{12.0, 6.0}
			checkArea(t, rectangle, 72.0)
		})
	t.Run(
		"Find circle area",
		func(t *testing.T) {
			circle := Circle{10}
			checkArea(t, circle, 314.1592653589793)
		})
}

func checkPerimeter(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	actual := shape.Perimeter()
	assertFloat64(t, expected, actual)
}

func checkArea(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	actual := shape.Area()
	assertFloat64(t, expected, actual)
}

func assertFloat64(t *testing.T, want float64, got float64) {
	if math.Abs(got-want) < math.SmallestNonzeroFloat64 {
		return
	}
	t.Errorf("want: %g got: %g", want, got)
}
