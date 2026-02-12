package main

import "fmt"

func ex1() {
	greetings := []string{"Hello", "Hola", "नमस्कार", "こんにちは", "Привіт"}
	first := greetings[:2]
	second := greetings[1:4]
	third := greetings[3:]
	fmt.Println("First", first)
	fmt.Println("Second", second)
	fmt.Println("Third", third)
}

func ex2() {
	message := "Hi 👩 and 👨"
	rs := []rune(message)
	fmt.Println(string(rs[3]))
}

type Employee struct {
	firstName string
	lastName  string
	id        int
}

func ex3() {
	one := Employee{}
	two := Employee{firstName: "ikradiua", lastName: "edian"}
	var three Employee
	three.firstName = "Test"
	three.lastName = "Test2"
	three.id = 10
	fmt.Println(one, two, three)
}

func main() {
	ex1()
	fmt.Println()
	ex2()
	fmt.Println()
	ex3()
}
