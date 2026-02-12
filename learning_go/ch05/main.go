package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	add = func(i, j int) int { return i + j }
	sub = func(i, j int) int { return i - j }
	mul = func(i, j int) int { return i * j }
	div = func(i, j int) (int, error) {
		var ret int
		if j == 0 {
			return ret, errors.New("division by zero")
		}
		return i / j, nil
	}
)

func ex01() {
	x := add(2, 3)
	fmt.Println(x)

	x, err := div(5, 0)
	if err != nil {
		fmt.Printf("got error %v\n", err)
	} else {
		fmt.Println(x)
	}
}

func fileLen(fname string) (flen int, err error) {
	fp, err := os.Open(fname)
	if err != nil {
		return flen, err
	}
	defer fp.Close()

	finfo, err := fp.Stat()
	if err != nil {
		return flen, err
	}

	return int(finfo.Size()), err
}

func ex02() {
	flen, err := fileLen("go.mod")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(flen)
	}
	flen, err = fileLen("go.modi")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(flen)
	}
}

func prefixer(prefix string) func(string) string {
	return func(input string) string {
		return prefix + " " + input
	}
}

func ex03() {
	helloPrefix := prefixer("Hello")
	fmt.Println(helloPrefix("Bob"))
	fmt.Println(helloPrefix("Maria"))
}

func main() {
	ex01()
	fmt.Println()

	ex02()
	fmt.Println()

	ex03()
	fmt.Println()
}
