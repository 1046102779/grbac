package slicelement

import (
	"math"
	"reflect"

	"github.com/pkg/errors"
)

type indexT struct{}

// int type
func (t *indexT) getIndexInt(data interface{}, element interface{}) (index int, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Int()
	for index = 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Int() == elem {
			return index, nil
		}
	}
	return -1, nil
}

// string type
func (t *indexT) getIndexString(data interface{}, element interface{}) (index int, err error) {
	dataVal := reflect.ValueOf(data)
	for index = 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.String() == element.(string) {
			return index, nil
		}
	}
	return -1, nil
}

// float32 type
func (t *indexT) getIndexFloat32(data interface{}, element interface{}) (index int, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Float()
	for index = 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if math.Abs(elemDataVal.Float()-elem) <= EPSINON {
			return index, nil
		}
	}
	return -1, nil
}

// uint type
func (t *indexT) getIndexUint(data interface{}, element interface{}) (index int, err error) {
	dataVal := reflect.ValueOf(data)
	elem := reflect.ValueOf(element).Uint()
	for index = 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Uint() == elem {
			return index, nil
		}
	}
	return -1, nil
}

// struct type
func (t *indexT) getIndexStruct(data interface{}, element interface{}, tag string) (index int, err error) {
	var (
		isExist bool = false
	)
	contain := new(contain)
	dataVal := reflect.ValueOf(data)
	for index = 0; index < dataVal.Len(); index++ {
		elemDataVal := reflect.Indirect(dataVal.Index(index))
		if elemDataVal.Kind() != reflect.Struct {
			err = errors.Wrap(errors.New("the underly type of data is not struct"), "isContainStruct")
			return
		}
		isExist, err = contain.decodeStruct(elemDataVal, element, tag)
		if err != nil {
			err = errors.Wrap(err, "getIndexStruct")
			return
		}
		if isExist {
			return index, nil
		}
	}
	return -1, nil
}

// decode struct type
func (t *indexT) decodeStruct(DataVal reflect.Value, element interface{}, tag string) (isExist bool, err error) {
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
