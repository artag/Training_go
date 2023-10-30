package shapes

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {
	perimeterTests := []struct {
		name     string
		shape    Shape
		expected float64
	}{
		{"Find rectangle perimeter", Rectangle{10.0, 5.0}, 30.0},
		{"Find circle perimeter", Circle{10.0}, 62.83185307179586},
	}

	for _, tt := range perimeterTests {
		t.Run(tt.name, func(t *testing.T) {
			checkPerimeter(t, tt.shape, tt.expected)
		})
	}
}

func TestArea(t *testing.T) {
	areaTests := []struct {
		name     string
		shape    Shape
		expected float64
	}{
		{"Find rectangle area", Rectangle{12.0, 6.0}, 72.0},
		{"Find circle area", Circle{10}, 314.1592653589793},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			checkArea(t, tt.shape, tt.expected)
		})
	}
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
