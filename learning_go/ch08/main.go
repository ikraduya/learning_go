package main

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

// Ex 1
//
//	type CanBeDoubled interface {
//		~int | ~int8 | ~int16 | ~int32 | ~int64 |
//			~uint | ~uint8 | ~uint16 | ~uint32 | ~uintptr |
//			~float32 | ~float64
//	}
type CanBeDoubled interface {
	constraints.Integer | constraints.Float
}

func double[T CanBeDoubled](v T) T {
	return v * 2
}

// Ex 2
type Printable interface {
	fmt.Stringer
	~int | ~float64
}

type PrintableInt int

func (i PrintableInt) String() string {
	return strconv.Itoa(int(i))
}

type PrintableFloat float64

func (f PrintableFloat) String() string {
	return strconv.FormatFloat(float64(f), 'f', 6, 64)
}

func PrintToScreen[T Printable](v T) {
	fmt.Println(v)
}

// Ex 3
// Singly linkedlist
type LinkedListData[T comparable] struct {
	Value       T
	NextPointer *LinkedListData[T]
}

type LinkedList[T comparable] struct {
	Head *LinkedListData[T]
}

func (l *LinkedList[T]) Add(v T) {
	if l.Head == nil {
		l.Head = &LinkedListData[T]{v, nil}
	} else {
		ptr := l.Head
		for ptr.NextPointer != nil {
			ptr = ptr.NextPointer
		}
		ptr.NextPointer = &LinkedListData[T]{v, nil}
	}
}

func (l *LinkedList[T]) Insert(v T, pos int) {

	var prevPtr *LinkedListData[T]
	ptr := l.Head
	var i int
	for i < pos && ptr.NextPointer != nil {
		i += 1
		prevPtr = ptr
		ptr = ptr.NextPointer
	}

	if i == pos {
		if prevPtr == nil { // at the start
			if ptr == nil {
				l.Head = &LinkedListData[T]{v, nil} // empty list
			} else {
				l.Head = &LinkedListData[T]{v, l.Head}
			}
			return
		}
		if ptr == nil { // at the end
			prevPtr.NextPointer = &LinkedListData[T]{v, nil}
			return
		}

		// in the middle
		nextPtr := prevPtr.NextPointer
		prevPtr.NextPointer = &LinkedListData[T]{v, nextPtr}
	} else if i < pos && ptr.NextPointer == nil {
		ptr.NextPointer = &LinkedListData[T]{v, nil}
	}
}

func (l *LinkedList[T]) Index(v T) int {
	if l.Head == nil {
		return -1
	}

	var i int
	ptr := l.Head
	for ptr != nil && ptr.Value != v {
		ptr = ptr.NextPointer
		i += 1
	}

	if ptr == nil {
		return -1
	}

	return i
}

func main() {
	fmt.Println(double(4))   // Example with int
	fmt.Println(double(4.5)) // Example with float64

	i := PrintableInt(1)
	f := PrintableFloat(6.234)
	PrintToScreen(i)
	PrintToScreen(f)

	// LinkedList sample usage
	var list LinkedList[int]

	// Add elements
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Insert at head (pos 0)
	list.Insert(5, 0)

	// Insert in the middle (pos 2)
	list.Insert(15, 2)

	list.Insert(40, 100)
	list.Insert(12, 5)
	list.Insert(123, 7)

	// Print list
	fmt.Print("list: ")
	for p := list.Head; p != nil; p = p.NextPointer {
		fmt.Printf("%d ", p.Value)
	}
	fmt.Println()

	// Index lookups
	fmt.Println("Index(15):", list.Index(15))
	fmt.Println("Index(999):", list.Index(999))
}
