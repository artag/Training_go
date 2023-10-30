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
		{name: "Find rectangle perimeter", shape: Rectangle{10.0, 5.0}, expected: 30.0},
		{name: "Find circle perimeter", shape: Circle{10.0}, expected: 62.83185307179586},
		{name: "Find triangle perimeter", shape: Triangle{3, 6, 7}, expected: 16},
	}

	for _, tt := range perimeterTests {
		t.Run(
			tt.name,
			func(t *testing.T) {
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
		{name: "Find rectangle area", shape: Rectangle{12.0, 6.0}, expected: 72.0},
		{name: "Find circle area", shape: Circle{10}, expected: 314.1592653589793},
		{name: "Find triangle area", shape: Triangle{3, 6, 7}, expected: 8.94427190999916},
	}

	for _, tt := range areaTests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				checkArea(t, tt.shape, tt.expected)
			})
	}
}

func checkPerimeter(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	actual := shape.Perimeter()
	assertFloat64(t, shape, expected, actual)
}

func checkArea(t *testing.T, shape Shape, expected float64) {
	t.Helper()
	actual := shape.Area()
	assertFloat64(t, shape, expected, actual)
}

func assertFloat64(t *testing.T, shape Shape, want, got float64) {
	if math.Abs(got-want) < math.SmallestNonzeroFloat64 {
		return
	}
	t.Errorf("Shape: %#v. want: %g got: %g", shape, want, got)
}
