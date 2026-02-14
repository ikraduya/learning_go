package main

/*
	#cgo LDFLAGS: -lm
	#include <string.h>
	#include <stdlib.h>

	int mini_calc(char *op, int a, int b) {
		if (strcmp(op, "+") == 0) {
			return a + b;
		}
		if (strcmp(op, "*") == 0) {
			return a * b;
		}
		if (strcmp(op, "-") == 0) {
			return a - b;
		}
		if (strcmp(op, "/") == 0) {
			if (b == 0) {
				return 0;
			}
			return a / b;
		}
		return 0;
	}
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

type StructA struct {
	StringMin6        string `minStrlen:"6"`
	StringMin2        string `minStrlen:"2"`
	StringNoMinLength string
	NotString         int
}

func ValidateStringLength(in any) error {
	var t reflect.Type = reflect.TypeOf(in)
	var tVal reflect.Value = reflect.ValueOf(in)
	var k reflect.Kind = tVal.Kind()
	if k != reflect.Struct {
		return nil
	}

	errs := make([]error, 0)

	for i := 0; i < tVal.NumField(); i++ {
		curFieldVal := tVal.Field(i)
		curFieldType := t.Field(i)
		if curFieldType.Type.Kind() != reflect.String {
			continue
		}

		minStrlenRaw, ok := curFieldType.Tag.Lookup("minStrlen")
		if !ok {
			continue
		}
		minStrlen, err := strconv.Atoi(minStrlenRaw)
		if err != nil {
			continue
		}

		value := curFieldVal.String()
		if len(value) < minStrlen {
			fieldName := curFieldType.Name
			errs = append(errs, fmt.Errorf("field %s with value `%s` has length %d is less than minStrLen=%d", fieldName, value, len(value), minStrlen))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func ex01() {
	a := StructA{
		"123s",
		"1",
		"",
		2,
	}
	err := ValidateStringLength(a)
	fmt.Println(err)
}

type OrderInfo struct {
	OrderCode   rune
	Amount      int
	OrderNumber uint16
	Items       []string
	IsReady     bool
}

type SmallOrderInfo struct {
	Items       []string
	Amount      int
	OrderCode   rune
	OrderNumber uint16
	IsReady     bool
}

func ex02() {
	fmt.Println(unsafe.Sizeof(OrderInfo{}))
	fmt.Println("OrderCode", unsafe.Offsetof(OrderInfo{}.OrderCode))
	fmt.Println("Amount", unsafe.Offsetof(OrderInfo{}.Amount))
	fmt.Println("OrderNumber", unsafe.Offsetof(OrderInfo{}.OrderNumber))
	fmt.Println("Items", unsafe.Offsetof(OrderInfo{}.Items))
	fmt.Println("IsReady", unsafe.Offsetof(OrderInfo{}.IsReady))

	fmt.Println(unsafe.Sizeof(SmallOrderInfo{}))
	fmt.Println("Items", unsafe.Offsetof(SmallOrderInfo{}.Items))
	fmt.Println("Amount", unsafe.Offsetof(SmallOrderInfo{}.Amount))
	fmt.Println("OrderCode", unsafe.Offsetof(SmallOrderInfo{}.OrderCode))
	fmt.Println("OrderNumber", unsafe.Offsetof(SmallOrderInfo{}.OrderNumber))
	fmt.Println("IsReady", unsafe.Offsetof(SmallOrderInfo{}.IsReady))
}

func ex03() {
	cString := C.CString("+")
	defer C.free(unsafe.Pointer(cString))

	ret := C.mini_calc(cString, 1, 2)
	fmt.Println("Returned", ret)
}

func main() {
	ex01()
	fmt.Println()

	ex02()
	fmt.Println()

	ex03()
	fmt.Println()

}
