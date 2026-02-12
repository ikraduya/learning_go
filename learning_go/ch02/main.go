package main

import (
	"fmt"
	"math"
)

func main() {
	// 1`
	var i int = 20
	var f float32 = 11.1
	i = int(f)
	fmt.Println(i, f)

	// 2
	const value = 50
	i = value
	f = value
	fmt.Println(i, f)

	// 3
	var b byte
	var smallI int32
	var bigI uint64
	b = math.MaxUint8
	smallI = math.MaxInt32
	bigI = math.MaxUint64
	b += 1
	smallI += 1
	bigI += 1
	fmt.Println(b, smallI, bigI)
}
