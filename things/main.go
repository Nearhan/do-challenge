package main

import (
	"fmt"
	"reflect"
)

type X2 struct {
	A []string
	B []string
}

type Outer struct {
	I map[string]X2
}

func main() {

	x := []string{"1", "2", "3"}
	y := []string{"1", "3", "2"}

	a := &X2{[]string{"1", "2", "3"}, nil}
	b := &X2{[]string{"1", "2", "3"}, nil}

	c := &Outer{map[string]X2{"dog": X2{[]string{"1", "2", "3"}, nil}}}
	d := &Outer{map[string]X2{"dog": X2{[]string{"1", "2", "3"}, nil}}}
	e := &Outer{map[string]X2{"dog": X2{[]string{"1", "2", "3"}, nil}, "cat": X2{}}}

	fmt.Println(reflect.DeepEqual(x, y))
	fmt.Println(reflect.DeepEqual(a, b))
	fmt.Println(reflect.DeepEqual(c, d))
	fmt.Println(reflect.DeepEqual(c, e))
}
