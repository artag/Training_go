# A Tour of Go

## Packages, variables, and functions

### Packages, Imports, Exported names

* Programs start running in package `main`.

* This code groups the imports into a parenthesized, "factored" import statement.

* Name is exported if it begins with a **capital** letter.

* "pi" do not start with a capital letter, so they are not exported.

* Any "unexported" names are not accessible from outside the package.

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
)

func init() {
    rand.Seed(71)
}

func main() {
    fmt.Println("My favorite number is", rand.Intn(10))
    fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))
    fmt.Println(math.Pi)
}
```

### Functions

* The type comes *after* the variable name.

* A function can return any number of results.

```go
func add(x int, y int) int {
    return x + y
}

func add2(x, y int) int {                       // shortened (x int, y int)
    return x + y
}

func swap(x, y string) (string, string) {       // return two results.
    return y, x
}

func main() {
    sum := add(42, 13)
    fmt.Println(sum)

    a, b := swap("hello", "world")
    fmt.Println(a, b)
}
```

### Named return values

* A return statement without arguments returns the named return values ("naked" return).

```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}

func main() {
    a, b := split(22)
    fmt.Printf("a = %d, b = %d", a, b)      // a = 9, b = 13
}
```

### Variables

* A var declaration can include initializers, one per variable.

* *Inside a function*, the `:=` short assignment statement can be used in place of a `var`
declaration with implicit type.

* *Outside a function*, every statement begins with a keyword (`var`, `func`, and so on)
and so the `:=` construct is not available.

```go
var c, python, java bool

var m, n int = 1, 2

func main() {
    i := 2
    fmt.Println(i, c, python, java)             // 2 false false false

    var golang, cpp, scala = true, false, "no!"
    fmt.Println(m, n, golang, cpp, scala)       // 1 2 true false no!
}
```

### Basic types

```text
bool

string

int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr

byte // alias for uint8

rune // alias for int32
     // represents a Unicode code point

float32 float64

complex64 complex128
```

* Variables declared without an explicit initial value are given their *zero value*.

The zero value is:

```text
0 for numeric types,
false for the boolean type, and
"" (the empty string) for strings.
```

```go
import "math/cmplx"

var (
    ToBe   bool       = false
    MaxInt uint64     = 1<<64 - 1
    z      complex128 = cmplx.Sqrt(-5 + 12i)
    p uintptr
)

func main() {
    fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)      // Type: bool Value: false
    fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)  // Type: uint64 Value: 18446744073709551615
    fmt.Printf("Type: %T Value: %v\n", z, z)            // Type: complex128 Value: (2+3i)
    fmt.Printf("Type: %T Value: %v\n", p, p)            // Type: uintptr Value: 0
}
```

### Type conversions

Some numeric conversions:

```text
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)
```

Or, put more simply:

```text
i := 42
f := float64(i)
u := uint(f)
```

### Type inference

```go
func main() {
    var i = 42          // int
    f := 3.142          // float64
    g := 0.867 + 0.5i   // complex128

    fmt.Printf("i is of type '%T'\n", i)    // i is of type 'int'
    fmt.Printf("f is of type '%T'\n", f)    // f is of type 'float64'
    fmt.Printf("g is of type '%T'\n", g)    // g is of type 'complex128'
}
```

### Constants

* Constants are declared like variables, but with the `const` keyword.

* Constants can be character, string, boolean, or numeric values.

* Constants cannot be declared using the `:=` syntax.

```go
const Pi = 3.14

// Еще можно задавать константы так
const (
    Big = 1 << 100
    Small = Big >> 99
)

func main() {
    const World = "мир"
    fmt.Println("Hello", World)
    fmt.Println("Happy", Pi, "Day")
}
```

## Flow control statements. `for`, `if`, `else`, `switch` and `defer`

## More types. `struct`, `slice`, and `map`

### Pointers

Go has no pointer arithmetic.

Pointer zero value is `nil`.

```go
func main() {
    i, j := 42, 2701

    p := &i                         // point to i
    fmt.Println("Read i:", *p)      // read i through the pointer
    *p = 21                         // set i through the pointer
    fmt.Println("Read new i:", i)   // see the new value of i

    p = &j                          // point to j
    *p = *p / 37                    // divide j through the pointer
    fmt.Println("Read new j:", j)   // see the new value of j
}
```

Вывод:

```text
Read i: 42
Read new i: 21
Read new j: 73
```

### Structs

* A `struct` is a collection of fields.

* Struct fields are accessed using a dot.

* Struct fields can be accessed through a struct pointer.

```go
type Vertex struct {
    X int
    Y int
}

func main() {
    v := Vertex{1, 2}
    fmt.Println(v)      // {1 2}

    v.X = 4
    fmt.Println(v.X)    // 4

    p := &v
    p.X = 10
    fmt.Println(v)      // {10 2}
    fmt.Println(*p)     // {10 2}
}
```

### Struct Literals

* You can list just a subset of fields by using the `Name:` syntax.

* The special prefix `&` returns a pointer to the struct value.

```go
type Vertex struct {
    X, Y int
}

var (
    v1 = Vertex{1, 2}       // has type Vertex
    v2 = Vertex{X: 1}       // Y:0 is implicit
    v3 = Vertex{}           // X:0 and Y:0
    p  = &Vertex{1, 2}      // has type *Vertex
)

func main() {
    fmt.Println(v1, p, v2, v3)      // {1 2} &{1 2} {1 0} {0 0}
}
```

### Arrays

* Arrays cannot be resized.

```go
func main() {
    var a [2]string
    a[0] = "Hello"
    a[1] = "World"
    fmt.Println(a[0], a[1])         // Hello World
    fmt.Println(a)                  // [Hello World]

    primes := [6]int{2, 3, 5, 7, 11, 13}
    evens := [6]int{2, 4, 6, 8}
    fmt.Println(primes)         // [2 3 5 7 11 13]
    fmt.Println(evens)          // [2 4 6 8 0 0]
}
```

### Slices

* Dynamically-sized

```go
func main() {
    primes := [6]int{2, 3, 5, 7, 11, 13}

    var s []int = primes[1:4]
    fmt.Println(s)      // [3 5 7]

    s2 := primes[0:2]
    fmt.Println(s2)     // [2 3]
}
```

### Slices are like references to arrays

* Changing the elements of a slice modifies the corresponding elements of its underlying array.

```go
func main() {
    names := [4]string{
        "John",
        "Paul",
        "George",
        "Ringo",
    }
    fmt.Println(names)      // [John Paul George Ringo]

    a := names[0:2]
    b := names[1:3]
    fmt.Println(a, b)       // [John Paul] [Paul George]

    b[0] = "XXX"
    fmt.Println(a, b)       // [John XXX] [XXX George]
    fmt.Println(names)      // [John XXX George Ringo]
}
```

### Slice literals

```go
func main() {
    q := []int{2, 3, 5, 7, 11, 13}
    fmt.Println(q)      // [2 3 5 7 11 13]

    r := []bool{true, false, true, true, false, true}
    fmt.Println(r)      // [true false true true false true]

    s := []struct {
        i int
        b bool
    }{
        {2, true},
        {3, false},
        {5, true},
        {7, true},
        {11, false},
        {13, true},
    }
    fmt.Println(s)      // [{2 true} {3 false} {5 true} {7 true} {11 false} {13 true}]
}
```

### Slice defaults

* Low bound - default value - 0
* High bound - default value - length of the slice

```go
func main() {
    s := []int{2, 3, 5, 7, 11, 13}
    a := s[:]
    fmt.Println(a)      // [2 3 5 7 11 13]

    b := s[1:4]
    fmt.Println(b)      // [3 5 7]

    c := s[:2]
    fmt.Println(c)      // [2 3]

    d := s[1:]
    fmt.Println(d)      // [3 5 7 11 13]
}
```

### Slice length and capacity

* The *length* of a slice is the number of elements it contains.

* The *capacity* of a slice is the number of elements in the underlying array,
counting from the first element in the slice.

```go
func main() {
    s := []int{2, 3, 5, 7, 11, 13}
    printSlice(s)           // len=6 cap=6 [2 3 5 7 11 13]

    // Slice the slice to give it zero length.
    s = s[:0]
    printSlice(s)           // len=0 cap=6 []

    // Extend its length.
    s = s[:4]
    printSlice(s)           // len=4 cap=6 [2 3 5 7]

    // Drop its first two values.
    s = s[2:]
    printSlice(s)           // len=2 cap=4 [5 7]
}

func printSlice(s []int) {
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
```

### Nil slices

* The zero value of a slice is `nil`.

* A `nil` slice has a length and capacity of 0 and has no underlying array.

```go
func main() {
    var s []int
    fmt.Println(s, len(s), cap(s))      // [] 0 0
    if s == nil {
        fmt.Println("nil!")             // nil!
    }
}
```

### Creating a slice with make

* The `make` function allocates a zeroed array and returns a slice that refers to that array.

```go
func main() {
    a := make([]int, 5)
    printSlice("a", a)      // a len=5 cap=5 [0 0 0 0 0]

    b := make([]int, 0, 5)
    printSlice("b", b)      // b len=0 cap=5 []

    c := b[:2]
    printSlice("c", c)      // c len=2 cap=5 [0 0]

    d := c[2:5]
    printSlice("d", d)      // d len=3 cap=3 [0 0 0]
}

func printSlice(s string, x []int) {
    fmt.Printf("%s len=%d cap=%d %v\n", s, len(x), cap(x), x)
}
```

### Slices of slices

```go
import "strings"

func main() {
    // Create a tic-tac-toe board.
    board := [][]string{
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
    }

    // The players take turns.
    board[0][0] = "X"           // X _ X
    board[2][2] = "O"           // O _ X
    board[1][2] = "X"           // _ _ O
    board[1][0] = "O"
    board[0][2] = "X"

    for i := 0; i < len(board); i++ {
        fmt.Printf("%s\n", strings.Join(board[i], " "))
    }
}
```

### Appending to a slice

```go
func main() {
    var s []int
    printSlice(s)       // len=0 cap=0 []

    // append works on nil slices.
    s = append(s, 0)
    printSlice(s)       // len=1 cap=1 [0]

    // The slice grows as needed.
    s = append(s, 1)
    printSlice(s)       // len=2 cap=2 [0 1]

    // We can add more than one element at a time.
    s = append(s, 2, 3, 4)
    printSlice(s)       // len=5 cap=6 [0 1 2 3 4]
}

func printSlice(s []int) {
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
```

### Range

* The `range` form of the `for` loop iterates over a slice or map.

* When ranging over a slice, two values are returned for each iteration.
The first is the *index*, and the second is a *copy of the element* at that index.

<table>
<tr>
<td>

```go
var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func main() {
    for i, v := range pow {
        fmt.Printf("[%d] 2^%d = %d\n", i, i, v)
    }
}
```

</td>
<td>

```text
[0] 2^0 = 1
[1] 2^1 = 2
[2] 2^2 = 4
[3] 2^3 = 8
[4] 2^4 = 16
[5] 2^5 = 32
[6] 2^6 = 64
[7] 2^7 = 128
```

</td>
</tr>
</table>

* You can skip the index or value by assigning to `_`

```text
for i, _ := range pow
for _, value := range pow
```

* If you only want the index, you can omit the second variable.

```text
for i := range pow
```

<table>
<tr>
<td>

```go
func main() {
    pow := make([]int, 10)
    for i := range pow {
        pow[i] = 1 << uint(i)       // == 2**i
    }
    for _, value := range pow {
        fmt.Printf("%d\n", value)
    }
}
```

</td>
<td>

```text
1
2
4
8
16
32
64
128
256
512
```

</td>
</tr>
</table>

### Maps

* A map maps keys to values.

* The zero value of a map is `nil`. A `nil` map has no keys, nor can keys be added.

* The `make` function returns a map of the given type, initialized and ready for use.

```go
type Vertex struct {
    Lat, Long float64
}

var m map[string]Vertex

func main() {
    m = make(map[string]Vertex)
    m["Bell Labs"] = Vertex{ 40.68433, -74.39967 }
    fmt.Println(m["Bell Labs"])     // {40.68433 -74.39967}
}
```

### Map literals

* If the top-level type is just a type name, you can omit it from the elements of the literal.

```go
type Vertex struct {
    Lat, Long float64
}

var m = map[string]Vertex {
    "Bell Labs": Vertex{ 40.68433, -74.39967 },
    "Google": Vertex{ 37.42202, -122.08408 },
}

// If the top-level type is just a type name, you can omit it from the elements of the literal.
var m = map[string]Vertex {
    "Bell Labs": {40.68433, -74.39967},
    "Google":    {37.42202, -122.08408},
}

func main() {
    fmt.Println(m)  // map[Bell Labs:{40.68433 -74.39967} Google:{37.42202 -122.08408}]
}
```

### Mutating Maps

* Insert or update an element in map m: `m[key] = elem`

* Retrieve an element: `elem = m[key]`

* Delete an element: `delete(m, key)`

* Test that a key is present with a two-value assignment: `elem, ok = m[key]`

  If `key` is in `m`, `ok` is `true`. If not, `ok` is `false`.

```go
func main() {
    m := make(map[string]int)

    m["Answer"] = 42
    fmt.Println("The value:", m["Answer"])      // The value: 42

    m["Answer"] = 48
    fmt.Println("The value:", m["Answer"])      // The value: 48

    delete(m, "Answer")
    fmt.Println("The value:", m["Answer"])      // The value: 0

    v, ok := m["Answer"]
    fmt.Println("The value:", v, "Present?", ok)    // The value: 0 Present? false
}
```

### Function values

* Function values may be used as function arguments and return values.

```go
import "math"

func compute(fn func(float64, float64) float64) float64 {
    return fn(3, 4)
}

func main() {
    hypot := func(x, y float64) float64 {
        return math.Sqrt(x * x + y * y)
    }
    fmt.Println(hypot(5, 12))       // 13   <= 5 * 5 + 12 * 12
    fmt.Println(compute(hypot))     // 5    <= 3 * 3 + 4 * 4
    fmt.Println(compute(math.Pow))  // 81   <= 3**4
}
```

### Function closures

Go functions may be closures. A closure is a function value that references variables
from outside its body. The function may access and assign to the referenced variables;
in this sense the function is "bound" to the variables.

The `adder` function returns a closure. Each closure is bound to its own `sum` variable.

<table>
<tr>
<td>

```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

func main() {
    pos, neg := adder(), adder()
    for i := 0; i < 10; i++ {
        fmt.Println(
            pos(i),
            neg(-2*i),
        )
    }
}
```

</td>
<td>

```text
0 0      // i = 0 -> sum = 0      | sum = 0
1 -2     // i = 1 -> sum = 0 + 1  | sum = 0 + -2
3 -6     // i = 2 -> sum = 1 + 2  | sum = -2 + -4
6 -12    // i = 3 -> sum = 3 + 3  | sum = -6 + -6
10 -20   // i = 4 -> sum = 6 + 4  | sum = -12 + -8
15 -30   // i = 5 -> sum = 10 + 5 | sum = -20 + -10
21 -42   // ...
28 -56
36 -72
45 -90
```

</td>
</tr>
</table>

## Methods and interfaces

### Methods

* Go does not have classes.

* A method is a function with a special *receiver* argument.

```go
import "math"

type MyFloat float64

type Vertex struct {
    X, Y float64
}

// Поведение будет аналогично функции: func Abs(v Vertex) float64
// Vertex - receiver argument
func (v Vertex) Abs() float64 {               
    return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

func (f MyFloat) AbsFloat() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())        // 5

    f := MyFloat(-math.Sqrt2)
    fmt.Println(f.AbsFloat())   // 1.4142135623730951
}
```

### Pointer receivers

* You can declare methods with pointer receivers (`*T`).

* Methods with pointer receivers can modify the value to which the receiver points.

* With a value receiver, the `Scale` method operates on a **copy** of the original `Vertex` value.

* Pointer receivers are more common than value receivers.

```go
type Vertex struct {
    X, Y float64
}

func (v *Vertex) Scale(f float64) {     // Can change Vertex
    v.X = v.X * f
    v.Y = v.Y * f
}

func (v Vertex) NotScale(f float64) {   // Operates on Vertex copy
    v.X = v.X * f
    v.Y = v.Y * f
}

func main() {
    v0 := Vertex{3, 4}
    var v1 = v0

    v0.Scale(10)
    v1.NotScale(10)

    fmt.Println("v0 = ", v0)    // v0 =  {30 40}
    fmt.Println("v1 = ", v1)    // v1 =  {3 4}
}
```

### Pointers and functions

```go
type Vertex struct {
    X, Y float64
}

func ScalePointer(v *Vertex, f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func ScaleCopy(v Vertex, f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func main() {
    v0 := Vertex{3, 4}
    v1 := Vertex{3, 4}

    ScalePointer(&v0, 10)       // Change Vertex
    ScaleCopy(v1, 10)

    fmt.Println("v0 = ", v0)    // v0 =  {30 40}
    fmt.Println("v1 = ", v1)    // v1 =  {3 4}
}
```

### Methods and pointer indirection

```go
type Vertex struct {
    X, Y float64
}

func (v *Vertex) Scale(f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func ScaleFunc(v *Vertex, f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func main() {
    v := Vertex{3, 4}
    v.Scale(2)
    ScaleFunc(&v, 10)

    p := &Vertex{4, 3}
    p.Scale(3)
    ScaleFunc(p, 10)

    fmt.Println(v, p)       // {60 80} &{120 90}
}
```

### Methods and pointer indirection (2)

```go
import "math"

type Vertex struct {
    X, Y float64
}

func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

func AbsFunc(v Vertex) float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())        // 5
    fmt.Println(AbsFunc(v))     // 5

    p := &Vertex{4, 3}
    fmt.Println(p.Abs())        // 5
    fmt.Println(AbsFunc(*p))    // 5
}
```

### Choosing a value or pointer receiver

* **Recommended** to use a **pointer** receiver.

  1) The method can modify the value
  2) Pointer can be more efficient than copying the value

* All methods on a given type should have either value or pointer receivers,
but not a mixture of both.

### Interfaces

* An *interface* type is defined as a set of method signatures.

* A value of interface type can hold any value that implements those methods.

```go
import "math"

type Abser interface {
    Abs() float64
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

type Vertex struct {
    X, Y float64
}

func (v *Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
    var a Abser
    f := MyFloat(-math.Sqrt2)
    v := Vertex{3, 4}

    a = f  // a MyFloat implements Abser
    fmt.Println(a.Abs())        // 1.4142135623730951

    a = &v // a *Vertex implements Abser
    fmt.Println(a.Abs())        // 5
}
```

### Interfaces are implemented implicitly

* There is no explicit declaration of intent, no "implements" keyword.

```go
type I interface {
    M()
}

type T struct {
    S string
}

// This method means type T implements the interface I,
// but we don't need to explicitly declare that it does so.
func (t T) M() {
    fmt.Println("Print: ", t.S)
}

func main() {
    var i I = T{"hello"}
    i.M()       // Print:  hello
}
```

## Tools

The main Go distribution includes tools for building, testing, and analyzing code:

* `go build`, which builds Go binaries using only information in the source files themselves, no separate makefiles
* `go test`, for unit testing and microbenchmarks
* `go fmt`, for formatting code
* `go install`, for retrieving and installing remote packages
* `go vet`, a static analyzer looking for potential errors in code
* `go run`, a shortcut for building and executing code
* `godoc`, for displaying documentation or serving it via HTTP
* `gorename`, for renaming variables, functions, and so on in a type-safe way
* `go generate`, a standard way to invoke code generators
* `go mod`, for creating a new module, adding dependencies, upgrading dependencies, etc.

It also includes *profiling* and *debugging* support, *runtime* instrumentation (for example,
to track *garbage collection* pauses), and a *race condition* tester.

Third-party tools (adds to the standard distribution):

* `gocode`, enables code autocompletion in many text editors
* `goimports`, automatically adds/removes package imports as needed
* `errcheck`, detects code that might unintentionally ignore errors.
