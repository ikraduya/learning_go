package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func MakePerson(firstName, lastName string, age int) Person {
	return Person{firstName, lastName, age}
}

func MakePersonPointer(firstName, lastName string, age int) *Person {
	return &Person{firstName, lastName, age}
}

func ex01() {
	p := MakePerson("abc", "def", 12)
	pp := MakePersonPointer("ikra", "duya", 69)
	fmt.Println(pp)
	fmt.Println(&p)
}

func UpdateSlice(s []string, item string) {
	s[len(s)-1] = item
	fmt.Println("UpdateSlice", s)
}

func GrowSlice(s []string, item string) {
	s = append(s, item)
	fmt.Println("GrowSlice", s)
}

func ex02() {
	a := []string{"1", "2", "3"}
	UpdateSlice(a, "bruh")
	fmt.Println(a)
	GrowSlice(a, "wat")
	fmt.Println(a) // kinda not work since the reference is different in GrowSlice
}

func ex03() {
	ps := make([]Person, 0, 10000000)
	for i := 0; i < 10000000; i++ {
		ps = append(ps, Person{FirstName: "ikra", LastName: "duya", Age: 16})
	}
}

func main() {
	ex01()
	fmt.Println()

	ex02()
	fmt.Println()

	ex03()
	fmt.Println()
}
