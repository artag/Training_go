package main

import "fmt"

type person struct {
	name string
	age  int
}

func newPerson(name string) *person {
	p := person{name: name}
	p.age = 42
	return &p
}

func main() {
	fmt.Println(person{"Bob", 20})
	fmt.Println(person{name: "Alice", age: 30})
	fmt.Println(person{name: "Fred"})
	fmt.Println(&person{name: "Ann", age: 40})
	fmt.Println(newPerson("John"))

	fmt.Println("---------")
	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	fmt.Println("---------")
	sp := &s
	fmt.Println("initial sp:", sp)
	fmt.Println("initial sp.age:", sp.age)

	sp.age = 51
	fmt.Println("changed s.age:", s.age)
	fmt.Println("changed sp.age:", sp.age)
}
