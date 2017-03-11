// Go library for finding element in slice type or operating set including union, interaction and difference.
//
// it not only supports the buildin types which includes []int/[]*int, []float/[]*float, []string/[]*string,
// but also it supports []struct/[]*struct .
// The latter is very important and convenient
package slicelement

import (
	"reflect"

	"github.com/pkg/errors"
)

// 校验输入数据的有效性
func checkInputValid(data interface{}, element interface{}) (err error) {
	// check data
	if err = checkData(data); err != nil {
		err = errors.Wrap(err, "checkDataValid")
		return
	}
	// check element
	if err = checkElement(element); err != nil {
		err = errors.Wrap(err, "checkDataValid")
		return
	}
	return
}

// check data
func checkData(data interface{}) (err error) {
	// check data: 1. 空指针
	if data == nil {
		err = errors.New("input: the first data isnil")
		return
	}
	// check data: 2. 数据类型必须为slice或者array类型
	dataVal := reflect.ValueOf(data)
	if dataVal.Kind() != reflect.Slice && dataVal.Kind() != reflect.Array {
		err = errors.New("input: type of the first data must be slice or array")
		return
	}
	return
}

// check element
func checkElement(element interface{}) (err error) {
	// if element is structure, it supports structure pointer
	// it doesn't support *int, *string, ...
	isPointer := false
	value := reflect.ValueOf(element)
	for value.Kind() == reflect.Ptr {
		if !isPointer {
			isPointer = true
		}
		if value.IsNil() {
			err = errors.New("the value of element is nil")
		}
		value = reflect.Indirect(value)
	}
	kind := getKindByValue(value)
	if kind != reflect.Int && kind != reflect.String &&
		kind != reflect.Float32 && kind != reflect.Uint {
		err = errors.Wrap(errors.New("element only supports `struct`, `int`, `string` and `float`"), "checkElement")
		return
	}
	return
}

func getKindByValue(val reflect.Value) (kind reflect.Kind) {
	kind = val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func getKindByKind(kind reflect.Kind) (dest reflect.Kind) {

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

// get the kind of underly data type
func getSliceUnderlyKind(data interface{}) (kind reflect.Kind, err error) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		err = errors.Wrap(errors.New("only support `slice` and `array`"), "getSliceUnderlyKind")
		return
	}
	typ := value.Type().Elem()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Kind(), nil
}

//  public method，the data whether exists element
//
// if data's type is a slice struct type, param `tag` can't be empty
func Contains(data interface{}, element interface{}, tag string) (isExist bool, err error) {
	// data only supports slice or array type
	if err = checkInputValid(data, element); err != nil {
		err = errors.Wrap(err, "Contains")
		return
	}
	contain := new(contain)
	kind, err := getSliceUnderlyKind(data)
	if err != nil {
		return false, err
	}
	kind = getKindByKind(kind)
	switch kind {
	case reflect.String:
		isExist, err = contain.isContainString(data, element)
	case reflect.Int:
		isExist, err = contain.isContainInt(data, element)
	case reflect.Uint:
		isExist, err = contain.isContainUint(data, element)
	case reflect.Float32:
		isExist, err = contain.isContainFloat32(data, element)
	case reflect.Struct:
		isExist, err = contain.isContainStructs(data, element, tag)
	}
	if err != nil {
		err = errors.Wrap(err, "Contains")
		return
	}
	return
}

// getting the element index in data, if not exists, return -1, nil
//
// if data's type is a slice struct type, param `tag` can't be empty
func GetIndex(data interface{}, element interface{}, tag string) (index int, err error) {
	// data only supports slice or array type
	if err = checkInputValid(data, element); err != nil {
		err = errors.Wrap(err, "Contains")
		return
	}
	indexInstance := new(indexT)
	kind, err := getSliceUnderlyKind(data)
	if err != nil {
		return -1, err
	}
	kind = getKindByKind(kind)
	switch kind {
	case reflect.String:
		index, err = indexInstance.getIndexString(data, element)
	case reflect.Int:
		index, err = indexInstance.getIndexInt(data, element)
	case reflect.Uint:
		index, err = indexInstance.getIndexUint(data, element)
	case reflect.Float32:
		index, err = indexInstance.getIndexFloat32(data, element)
	case reflect.Struct:
		index, err = indexInstance.getIndexStruct(data, element, tag)
	}
	if err != nil {
		err = errors.Wrap(err, "Contains")
		return
	}
	return
}
