package slicelement

import (
	"fmt"
)

func ExampleGetIndex_int() {
	var (
		data []int = []int{1, 2, 3, 4, 5}
		elem int   = 2
	)
	index, err := GetIndex(data, elem, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("index=%d\n", index)
	// output: index=1
	return
}

func ExampleGetIndex_uint() {
	var (
		data []uint = []uint{1, 2, 3, 4, 5}
		elem uint   = 2
	)
	index, err := GetIndex(data, elem, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("index=%d\n", index)
	// output: index=1
	return
}

func ExampleGetIndex_string() {
	var (
		data []string = []string{"abc", "def", "hig"}
		elem string   = "def"
	)
	index, err := GetIndex(data, elem, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("index=%d\n", index)
	// output: index=1
	return
}

func ExampleGetIndex_float32() {
	var (
		data []float32 = []float32{1, 2, 3, 4, 5}
		elem float32   = 2
	)
	index, err := GetIndex(data, elem, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("index=%d\n", index)
	// output: index=1
	return
}

func ExampleGetIndex_struct() {
	type Person struct {
		Name     string
		Age      int
		Children []string
	}
	data := []*Person{
		&Person{
			Name:     "John",
			Age:      29,
			Children: []string{"David", "Lily", "Bruce Lee"},
		},
		&Person{
			Name:     "Joe",
			Age:      18,
			Children: []string{},
		},
	}
	elem := 18
	index, err := GetIndex(data, elem, "Age")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("index=%d\n", index)
	// output: index=1
}
