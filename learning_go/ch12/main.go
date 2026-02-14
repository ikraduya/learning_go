package main

import (
	"context"
	"fmt"
	"math"
	"sync"
)

func ex01() {
	var writers sync.WaitGroup
	var reader sync.WaitGroup
	ch := make(chan int)
	writer := func(start int) {
		defer writers.Done()
		for i := start; i < start+10; i++ {
			ch <- i
		}
	}

	writers.Add(2)
	go writer(0)   // first goroutine
	go writer(100) // second goroutine

	reader.Add(1)
	go func() {
		defer reader.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}() // third goroutine

	writers.Wait()
	close(ch)
	reader.Wait()
}

func ex02() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	var writersWG sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	writer := func(ch *chan int, start int) {
		defer close(*ch)
		defer writersWG.Done()
		for i := start; i < start+10; i++ {
			*ch <- i
		}
	}
	writersWG.Add(2)
	go writer(&ch1, 0)
	go writer(&ch2, 100)

	go func() {
		defer cancel()
		writersWG.Wait()
	}()

reader:
	for {
		select {
		case v, ok := <-ch1:
			if ok {
				fmt.Println("From writer 1", v)
			} else {
				ch1 = nil
			}
		case v, ok := <-ch2:
			if ok {
				fmt.Println("From writer 2", v)
			} else {
				ch2 = nil
			}
		case <-ctx.Done():
			break reader
		}
	}
}

func buildDict() map[int]float64 {
	dict := make(map[int]float64)
	for i := 0; i < 100_000; i++ {
		dict[i] = math.Sqrt(float64(i))
	}

	return dict
}

var buildDictCached func() map[int]float64 = sync.OnceValue(buildDict)

func ex03() {
	dict := buildDictCached()
	for i := 0; i < 100_000; i += 1000 {
		fmt.Println(dict[i])
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
