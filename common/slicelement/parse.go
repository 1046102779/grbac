package slicelement

import (
	"math"
	"reflect"

	"github.com/pkg/errors"
)

type contain struct{}

const (
	EPSINON float64 = 0.00001
)

// string type
func (t *contain) isContainString(data interface{}, element interface{}) (isExist bool, err error) {
	dataVal := reflect.ValueOf(data)
	for index := 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.String() == element.(string) {
			return true, nil
		}
	}
	return false, nil
}

// int type
func (t *contain) isContainInt(data interface{}, element interface{}) (isExist bool, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Int()
	for index := 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Int() == elem {
			return true, nil
		}
	}
	return false, nil
}

// float32 type
func (t *contain) isContainFloat32(data interface{}, element interface{}) (isExist bool, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Float()
	for index := 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if math.Abs(elemDataVal.Float()-elem) <= EPSINON {
			return true, nil
		}
	}
	return false, nil
}

// uint type
func (t *contain) isContainUint(data interface{}, element interface{}) (isExist bool, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Uint()
	for index := 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Uint() == elem {
			return true, nil
		}
	}
	return false, nil
}

// struct type
func (t *contain) isContainStructs(data interface{}, element interface{}, tag string) (isExist bool, err error) {
	dataVal := reflect.ValueOf(data)
	for index := 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Kind() != reflect.Struct {
			err = errors.Wrap(errors.New("the underly type of data is not struct"), "isContainStruct")
			return
		}
		isExist, err = t.decodeStruct(elemDataVal, element, tag)
		if err != nil {
			err = errors.Wrap(err, "isContainStructs")
			return
		}
		if isExist {
			return true, nil
		}
	}
	return false, nil
}

// decode struct type
func (t *contain) decodeStruct(DataVal reflect.Value, element interface{}, tag string) (isExist bool, err error) {
	noTag := false
	if DataVal.Kind() != reflect.Struct {
		err = errors.Wrap(errors.New("the value's kind is not struct"), "decodeStruct")
		return
	}
	target := reflect.ValueOf(element)
	tagKind := getKindByValue(target)
	for index := 0; index < DataVal.NumField(); index++ {
		dataElemVal := DataVal.Field(index)
		typ := DataVal.Type()
		if typ.Field(index).Name == tag {
			noTag = true
			switch kind := getKindByValue(dataElemVal); kind {
			case reflect.String:
				if tagKind == reflect.String && target.String() == dataElemVal.String() {
					return true, nil
				} else if tagKind != reflect.String {
					err = errors.Wrap(errors.New("`"+tag+"` type is not string"), "decodeStruct")
					return false, err
				}
			case reflect.Int:
				if tagKind == reflect.Int && target.Int() == dataElemVal.Int() {
					return true, nil
				} else if tagKind != reflect.Int {
					err = errors.Wrap(errors.New("`"+tag+"` type is not int"), "decodeStruct")
					return false, err
				}
			case reflect.Uint:
				if tagKind == reflect.Uint && target.Uint() == dataElemVal.Uint() {
					return true, nil
				} else if tagKind != reflect.Uint {
					err = errors.Wrap(errors.New("`"+tag+"` type is not uint"), "decodeStruct")
					return false, err
				}
			case reflect.Float32:
				if tagKind == reflect.Float32 && math.Abs(target.Float()-dataElemVal.Float()) < EPSINON {
					return true, nil
				} else if tagKind != reflect.Float32 {
					err = errors.Wrap(errors.New("`"+tag+"` type is not float"), "decodeStruct")
					return false, nil
				}
			}
		}
	}
	if !noTag {
		err = errors.Wrap(errors.New("no exist in the tag `"+tag+"`"), "decodeStruct")
		return
	}
	return false, nil
}
