# Integers

## Запуск примеров из тестов

Можно привести примеры в качестве документации.

Часто примеры кода можно найти за пределами кодовой базы.
Примеры компилируются (и при необходимости выполняются) как часть
тестов пакета.

Также как и тесты, примеры это функции, которые находятся в
файлах `*_test.go`.

Пример example (в тестах):

```go
func ExampleAdd() {
    sum := Add(1, 5)
    fmt.Println(sum)
    // Output: 6
}
```

Запустим тесты и примеры:

```text
go test -v
```

Пример запустится в составе тестов, если будет присутствовать комментарий `// Output: 6`.
Если этот комментарий не указать, то этот пример будет компилироваться, но не будет запускаться.

Сравнение результатов функции идет со стандартным выводом.
Сравнение игнорирует пробелы в начале и в конце.

### Еще examples

```go
func ExampleSalutations() {
    fmt.Println("hello, and")
    fmt.Println("goodbye")
    // Output:
    // hello, and
    // goodbye
}
```

Unordered output - допускается любой порядок строк.

```go
func ExamplePerm() {
    for _, value := range Perm(5) {
        fmt.Println(value)
    }
    // Unordered output: 4
    // 2
    // 1
    // 3
    // 0
}
```

### Наименование examples

Соглашения по наименованию examples для package.
`F` - функция, `T` - тип, метод `M` типа `T`:

```go
func Example() { ... }
func ExampleF() { ... }
func ExampleT() { ... }
func ExampleT_M() { ... }
```

Несколько примеров функций для пакета/типа/функции/метода могут быть предоставлены путем добавления к имени отдельного суффикса. Суффикс должен начинаться со строчной буквы:

```go
func Example_suffix() { ... }
func ExampleF_suffix() { ... }
func ExampleT_suffix() { ... }
func ExampleT_M_suffix() { ... }
```

## Создание документации

Комментарии над функциями позволяет их документировать.

Запуск:

```text
godoc -http=:8000
```

В браузере:

```text
http://localhost:8000/pkg/
```

Можно увидеть сгенеренную документацию по пакету `integers`.

## Wrapping up

* Integers, addition
* Writing better documentation so users of our code can understand its usage quickly
* Examples of how to use our code, which are checked as part of our tests
