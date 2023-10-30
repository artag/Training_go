package shapes

import "math"

type Shape interface {
	Perimeter() float64
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Triangle struct {
	sideA float64
	sideB float64
	sideC float64
}

func (t Triangle) Perimeter() float64 {
	return t.sideA + t.sideB + t.sideC
}

func (t Triangle) Area() float64 {
	var s = (t.sideA + t.sideB + t.sideC) * 0.5
	var area = math.Sqrt(s * (s - t.sideA) * (s - t.sideB) * (s - t.sideC))
	return area
}
