# Arrays and slices

Arrays allow you to store multiple elements of the same type in a
variable in a particular order.

## Инициализация массивов

Arrays have a *fixed capacity* which you define when you declare the
variable. We can initialize an array in two ways:

```go
numbers := [5]int{1, 2, 3, 4, 5}

numbers := [...]int{1, 2, 3, 4, 5}
```

## Arrays print format

Using the `%v` placeholder to print the "default" format, which works well for arrays.

Вывод массива будет наподобие: `[1 2 3 4 5]`.

## range

`range` lets you iterate over an array. On each iteration, range returns
two values - the *index* and the *value*.

Примеры:

```go
func main() {
    nums := []int{2, 3, 4}
    sum := 0
    for _, num := range nums {
        sum += num
    }
    fmt.Println("sum:", sum)

    for i, num := range nums {
        if num == 3 {
            fmt.Println("index:", i)
        }
    }

    kvs := map[string]string{"a": "apple", "b": "banana"}
    for k, v := range kvs {
        fmt.Printf("%s -> %s\n", k, v)
    }

    for k := range kvs {
        fmt.Println("key:", k)
    }

    for i, c := range "go" {
        fmt.Println(i, c)
    }
}
```

## Slice

Slice allows us to have collections of any size.

## Инициализация slice

```go
mySlice := []int{1,2,3}
```

Или с помощью `make`:

```go
slice := make([]int, capacity)
```

`capacity` - начальный размер slice.

## `append`

Функция `append` берет два slice и возвращает новый объединенный slice,
содержащий все элементы двух предыдущих slice'ов.

## Покрытие тестами. Coverage tool

Запуск тестов с измерением покрытия:

```text
go test -cover
```

## Variadic functions

**Variadic functions** can be called with any number of trailing arguments.
For example, `fmt.Println` is a common variadic function.

```go
func sum(nums ...int) {
    fmt.Print(nums, " ")
    total := 0

    for _, num := range nums {
        total += num
    }
    fmt.Println(total)
}

func main() {
    sum(1, 2)                       // [1 2] 3
    sum(1, 2, 3)                    // [1 2 3] 6

    nums := []int{1, 2, 3, 4}       // [1 2 3 4] 10
    sum(nums...)
}
```

## Как сравнить arrays and slices. `func DeepEqual` (`reflect``)

Сделать итерацию по массиву/срезу с сравнением каждого элемента. Или воспользоваться функцией
из `reflect`:

```go
func DeepEqual(x, y any) bool
```

Данную функцию удобно использовать в тестах.

DeepEqual reports whether x and y are “deeply equal,” defined as follows.
Two values of identical type are deeply equal if one of the following cases applies.
Values of distinct types are never deeply equal.

Типы данных в `DeepEqual` не проверяются компилятором, только во время выполнения.

## Дополнительно

* [Go Slices: usage and internals](https://go.dev/blog/slices-intro)
