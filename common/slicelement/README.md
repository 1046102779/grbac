#  package slicelement
Go library for finding element in slice type or operating set including union, interaction and difference.

it not only supports the buildin types which includes `[]int/[]*int`, `[]float/[]*float`, `[]string/[]*string`, but also it supports `[]struct/[]*struct` . The latter is very important and convenient
## Installation

Standard  `go get`:

```
    $  go get -v -u github.com/1046102779/slicelement
```

## Index

```go
//  find the element whether exists in data, if exist, return true, nil 
func Contains(data interface{}, elem interface{}, tag string) (bool, error)

// get the element index in data, if not exist, return -1, nil. 
func GetIndex(data interface{}, elem interface{}, tag string) (int, error)

// set difference, formula:  dataC =  dataA - dataB
// Param tag: unique key.  if not , it may be covered.
// eg. type User struct { UserId int, Name string, Tel string, Sex int16}
// var ( usersA, usersB []*User )
func GetDifference(dataA interface{}, dataB interface{}, tag string) (interface{}, error)

// set interaction, formula: dataC = dataA ∩ dataB , it also supports slice struct
// Param tag: unique key.  if not , it may be covered.
// eg. type User struct { UserId int, Name string, Tel string, Sex int16}
// var ( usersA, usersB []*User )
func GetInteraction(dataA interface{}, dataB interface{}, tagName string) (interface{}, error) 

// set union, formula: dataC = dataA U dataB
// Param tag: unique key.  if not , it may be covered.
// eg. type User struct { UserId int, Name string, Tel string, Sex int16}
// var ( usersA, usersB []*User )
func GetUnion(dataA interface{}, dataB interface{}, tagName string) ( interface{}, error)

desc: if the data's type is not []*struct/[]struct, the tag value is empty
```


## Usage & Example

For usage and examples see the [Godoc](https://godoc.org/github.com/1046102779/slicelement)

###  example 1:  []struct  GetIndex
```go
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
index, err := slicelement.GetIndex(data, elem, "Age")
if err != nil {
    fmt.Println(err.Error())
}
fmt.Println("index=", index)
// output: index=1
```

###  example 2:  []struct GetUnion 
```go
type Student struct {
    Name string
    Age  int
}
studentA := []Student{             studentB := []Student{
    Student{                            Student{
        Name: "donghai",                    Name: "Joe",
        Age:  29,                           Age:  18,
    },                                  },
    Student{                            Student{
        Name: "jixaing",                    Name: "David",
        Age:  19,                           Age:  19,
    },                                  },
}                                   }

if studentC, err := slicelement.GetUnion(studentA, studentB, "Age"); err != nil {
    fmt.Println(err.Error())
} else {
    fmt.Println("result: ", studentC)
}
// result:  [{donghai 29} {jixaing 19} {Joe 18}]  
// {"David", 19} is covered by {"jixiang", 19}
```
