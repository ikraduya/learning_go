package main

import (
	"fmt"
	"math/rand"
)

func ex01() []int {
	s := make([]int, 0, 100)
	for i := 0; i < cap(s); i++ {
		s = append(s, rand.Int()%101)
	}
	fmt.Println(s)

	return s
}

func ex02(s []int) {
loop:
	for _, v := range s {
		switch {
		case v%2 == 0 && v%3 == 0:
			fmt.Println("Six!")
			break loop
		case v%2 == 0:
			fmt.Println("Two!")
		case v%3 == 0:
			fmt.Println("Three!")
		default:
			fmt.Println("Never mind")
		}
	}
}

func ex03() {
	var total int
	for i := 0; i < 10; i++ {
		total := total + 1 // bug due to shadowing
		fmt.Println(total)
	}
}

func main() {
	s := ex01()
	ex02(s)
	ex03()
}
