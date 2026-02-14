//go:build !test

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// set everything up
	ch1 := make(chan []byte, 100)
	ch2 := make(chan Result, 100)
	go DataProcessor(ch1, ch2)
	f, err := os.Create("results.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	go WriteData(ch2, f)
	err = http.ListenAndServe(":8080", NewController(ch1))
	if err != nil {
		fmt.Println(err)
	}
}
