# Structs, methods & interfaces

## Struct

A **struct** is just a named collection of fields where you can store data.

Пример:

```go
// An empty struct.
struct {}

// A struct with 6 fields.
struct {
    x, y int
    u float32
    _ float32  // padding
    A *[]int
    F func()
}
```

## Methods

A **method** - is a function with a receiver.

c

```go
// Struct
type Circle struct {
    Radius float64
}

// Cirle methodfunc (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

## Interfaces

An **interface** type defines a *type set*.

Пример:

```go
// Interface
type Shape interface {
    Area() float64
}

// Struct. Satisfies the Shape interface (has method Area)
type Circle struct {
    Radius float64
}

// Cirle method
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

## Table driven tests

Пример:

```go
func TestArea(t *testing.T) {
    areaTests := []struct {
        name string     // test name
        shape Shape     // shape
        want float64    // expected value
    }{
        {name: "rectangle test", shape: Rectangle{12, 6}, want: 72.0},
        {name: "circle test", shape: Circle{10}, want: 314.1592653589793},
    }

    for _, tt := range areaTests {
        t.Run(
            tt.name,
            func(t *testing.T) {
                got := tt.shape.Area()
                if got != tt.want {
                    t.Errorf("got %g want %g", got, tt.want)
                }
            })
    }
}
```

## `fmt` options

- `%f` - print float
- `%.2f` - print float with 2 decimal places
- `%g` - print a more precise decimal number than `f`
- `%#v` - print out struct with the values in its field
