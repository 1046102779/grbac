package slicelement

import (
	"fmt"
)

// union
func ExampleGetUnion_int() {
	dataA := []int{1, 2, 3, 4, 5}
	dataB := []int{2, 4, 6, 7}
	temp, err := GetUnion(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [1 2 3 4 5 6 7]
	return
}

func ExampleGetUnion_uint() {
	dataA := []uint{1, 2, 3, 4, 5}
	dataB := []uint{2, 4, 6, 7}
	temp, err := GetUnion(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [1 2 3 4 5 6 7]
	return
}

func ExampleGetUnion_float32() {
	dataA := []float32{1, 2, 3, 4, 5}
	dataB := []float32{2, 4, 6, 7}
	temp, err := GetUnion(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [1 2 3 4 5 6 7]
	return
}

func ExampleGetUnion_string() {
	str1, str2, str3, str4, str5 := "1", "2", "3", "4", "5"
	dataA := []*string{&str1, &str2, &str3}
	dataB := []*string{&str2, &str3, &str4, &str5}
	temp, err := GetUnion(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	value := temp.([]*string)
	for index := 0; index < len(value); index++ {
		fmt.Printf("%s ", *value[index])
	}
	// output: 1 2 3 4 5
}

func ExampleGetUnion_struct() {
	type Student struct {
		Name string
		Age  int
	}
	studentA := []Student{
		Student{
			Name: "donghai",
			Age:  29,
		},
		Student{
			Name: "jixaing",
			Age:  19,
		},
	}

	studentB := []Student{
		Student{
			Name: "Joe",
			Age:  18,
		},
		Student{
			Name: "David",
			Age:  19,
		},
	}
	if value, err := GetUnion(studentA, studentB, "Age"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("studentA U studentB, result:", value)
	}
	// output: studentA U studentB, result: [{donghai 29} {jixaing 19} {Joe 18}]
}

// interaction
func ExampleGetInteraction_int() {
	dataA := []int{1, 2, 3, 4, 5}
	dataB := []int{2, 4, 6, 7}
	temp, err := GetInteraction(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [2 4]
	return
}

func ExampleGetInteraction_uint() {
	dataA := []uint{1, 2, 3, 4, 5}
	dataB := []uint{2, 4, 6, 7}
	temp, err := GetInteraction(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [2 4]
	return
}

func ExampleGetInteraction_float32() {
	dataA := []float32{1, 2, 3, 4, 5}
	dataB := []float32{2, 4, 6, 7}
	temp, err := GetInteraction(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("result:", temp)
	// output: result: [2 4]
	return
}

func ExampleGetInteraction_string() {
	str1, str2, str3, str4, str5 := "1", "2", "3", "4", "5"
	dataA := []*string{&str1, &str2, &str3}
	dataB := []*string{&str2, &str3, &str4, &str5}
	temp, err := GetInteraction(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("temp: ", temp)
}

func ExampleGetInteraction_struct() {
	type Student struct {
		Name string
		Age  int
	}
	studentA := []Student{
		Student{
			Name: "donghai",
			Age:  29,
		},
		Student{
			Name: "jixaing",
			Age:  19,
		},
	}

	studentB := []Student{
		Student{
			Name: "Joe",
			Age:  18,
		},
		Student{
			Name: "David",
			Age:  19,
		},
	}
	if value, err := GetInteraction(studentA, studentB, "Age"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("studentA U studentB, result: ", value)
	}
}

// difference
func ExampleGetDifference_int() {
	dataA := []int{1, 2, 3, 4, 5}
	dataB := []int{2, 4, 6, 7}
	temp, err := GetDifference(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("temp: ", temp)
	return
}

func ExampleGetDifference_uint() {
	dataA := []uint{1, 2, 3, 4, 5}
	dataB := []uint{2, 4, 6, 7}
	temp, err := GetDifference(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("temp: ", temp)
	return
}

func ExampleGetDifference_float32() {
	dataA := []float32{1, 2, 3, 4, 5}
	dataB := []float32{2, 4, 6, 7}
	temp, err := GetDifference(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("temp: ", temp)
	return
}

func ExampleGetDifference_string() {
	str1, str2, str3, str4, str5 := "1", "2", "3", "4", "5"
	dataA := []*string{&str1, &str2, &str3}
	dataB := []*string{&str2, &str3, &str4, &str5}
	temp, err := GetDifference(dataA, dataB, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("temp: ", temp)
	return
}

func ExampleGetDifference_struct() {
	type Student struct {
		Name string
		Age  int
	}
	studentA := []Student{
		Student{
			Name: "donghai",
			Age:  29,
		},
		Student{
			Name: "jixaing",
			Age:  19,
		},
	}

	studentB := []Student{
		Student{
			Name: "Joe",
			Age:  18,
		},
		Student{
			Name: "David",
			Age:  19,
		},
	}
	if value, err := GetDifference(studentA, studentB, "Age"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("studentA U studentB, result: ", value)
	}
}
