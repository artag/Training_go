# Iterations

Примеры разных форм использования `for`:

```go
func main() {
    i := 1
    for i <= 3 {
        fmt.Println(i)
        i = i + 1
    }

    for j := 7; j <= 9; j++ {
        fmt.Println(j)
    }

    for {
        fmt.Println("loop")
        break
    }

    for n := 0; n <= 5; n++ {
        if n%2 == 0 {
            continue
        }
        fmt.Println(n)
    }
}
```

## Создание и запуск benchmark

Пишется в тестах.

Функция должна быть вида:

```go
func BenchmarkXxx(*testing.B)
```

Benchmark запускаются последовательно.

Пример benchmark:

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

Запуск benchmark `go test -bench=.`.

На Windows, в Powershell `go test -bench="."`.

### Дополнительно

Сброс таймера внутри benchmark. Например, если требуется долгая инициализация перед
измеряемым методом. Пример:

```go
func BenchmarkBigLen(b *testing.B) {
    big := NewBig()
    b.ResetTimer()      // Сброс таймера
    for i := 0; i < b.N; i++ {
        big.Len()
    }
}
```

Параллельный benchmark можно сделать через вспомогательную функцию `RunParallel`.
Такие тесты запускаюся через флаг `go test -cpu`. Пример:

```go
func BenchmarkTemplateParallel(b *testing.B) {
    templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
    b.RunParallel(func(pb *testing.PB) {
        var buf bytes.Buffer
        for pb.Next() {
            buf.Reset()
            templ.Execute(&buf, "World")
        }
    })
}
```
