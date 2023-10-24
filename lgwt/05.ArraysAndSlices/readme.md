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

## Покрытие тестами. Coverage tool

Запуск тестов с измерением покрытия:

```text
go test -cover
```
