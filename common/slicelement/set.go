package slicelement

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// check input dataA, if data is nil, it will new object
func checkSetDataA(data interface{}) (err error) {
	value := reflect.ValueOf(data)
	if data == nil {
		// it need new object
		underlyType := value.Type().Elem()
		newValue := reflect.MakeSlice(underlyType, 0, 0)
		if value.CanSet() {
			value.Set(newValue)
		}
		return
	}
	if value.Kind() != reflect.Slice {
		err = errors.New("the first input data must be slice type")
		return
	}
	return
}

// check input dataB, if dataB is nil and check whether the data  need to be added
func checkSetDataB(data interface{}) (needAdd bool, err error) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		err = errors.New("the second data must be slice or array")
		return
	}
	if value.IsNil() {
		return false, nil
	}
	return true, nil
}

// check two input datas
func checkSetInputData(dataA interface{}, dataB interface{}) (needAdd bool, err error) {
	// check dataA
	if err = checkSetDataA(dataA); err != nil {
		err = errors.Wrap(err, "checkSetInputData")
		return
	}
	// check dataB
	if needAdd, err = checkSetDataB(dataB); err != nil {
		err = errors.Wrap(err, "checkSetInputData")
		return false, err
	} else if !needAdd {
		return false, nil
	}
	// if dataA and dataB is not the same, return err
	underlyTypeA, err := getSliceUnderlyKind(dataA)
	if err != nil {
		err = errors.Wrap(err, "checkSetInputData")
		return
	}
	underlyTypeB, err := getSliceUnderlyKind(dataB)
	if err != nil {
		err = errors.Wrap(err, "checkSetInputData")
		return
	}
	if underlyTypeA != underlyTypeB {
		err = errors.New("input datas type are not the same")
		return
	}
	return
}

// get interaction . formula result = dataA U dataB
// both dataA and dataB are a slice type
func GetUnion(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	var needAdd bool = false
	if needAdd, err = checkSetInputData(dataA, dataB); err != nil {
		err = errors.Wrap(err, "GetUnion")
		return
	} else if !needAdd {
		return dataA, nil
	}
	union := &union{}
	kindA, err := getSliceUnderlyKind(dataA)
	if err != nil {
		err = errors.Wrap(err, "GetUnion")
		return
	}
	kindA = getKindByKind(kindA)
	switch kindA {
	case reflect.Uint, reflect.Int, reflect.Float32, reflect.String:
		result, err = union.getNonStruct(dataA, dataB)
	case reflect.Struct:
		result, err = union.getStruct(dataA, dataB, tagName)
	}
	return
}

// get interaction . formula result = dataA N dataB
// both dataA and dataB are a slice type
func GetInteraction(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	var needAdd bool = false
	if needAdd, err = checkSetInputData(dataA, dataB); err != nil {
		err = errors.Wrap(err, "Union")
		return
	} else if !needAdd {
		return dataA, nil
	}
	interaction := &interaction{}
	kindA, err := getSliceUnderlyKind(dataA)
	if err != nil {
		err = errors.Wrap(err, "GetInteraction")
		return
	}
	kindA = getKindByKind(kindA)
	switch kindA {
	case reflect.Uint, reflect.Int, reflect.Float32, reflect.String:
		result, err = interaction.getNonStruct(dataA, dataB)
	case reflect.Struct:
		result, err = interaction.getStruct(dataA, dataB, tagName)
	}
	return
}

// get difference . formula result = dataA - dataB
// both dataA and dataB are a slice type
func GetDifference(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	var needAdd bool = false
	if needAdd, err = checkSetInputData(dataA, dataB); err != nil {
		err = errors.Wrap(err, "GetDifference")
		return
	} else if !needAdd {
		return dataA, nil
	}
	diff := &difference{}
	kindA, err := getSliceUnderlyKind(dataA)
	if err != nil {
		err = errors.Wrap(err, "GetDifference")
		return
	}
	kindA = getKindByKind(kindA)
	switch kindA {
	case reflect.Int, reflect.Float32, reflect.String:
		result, err = diff.getNonStruct(dataA, dataB)
	case reflect.Struct:
		result, err = diff.getStruct(dataA, dataB, tagName)
	}
	return
}

// union
type union struct{}

func (t *union) getNonStruct(dataA interface{}, dataB interface{}) (result interface{}, err error) {
	valueB := reflect.ValueOf(dataB)
	resultVal := reflect.ValueOf(dataA)
	for index := 0; index < valueB.Len(); index++ {
		var subIndex int = 0
		for subIndex = 0; subIndex < resultVal.Len(); subIndex++ {
			resultElem := reflect.Indirect(resultVal.Index(subIndex))
			bElem := reflect.Indirect(valueB.Index(index))
			if resultElem.Interface() == bElem.Interface() {
				break
			}
		}
		if subIndex == resultVal.Len() {
			resultVal = reflect.Append(resultVal, valueB.Index(index))
		}
	}
	return resultVal.Interface(), nil
}

// tagName: unique key
func (t *union) getStruct(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	if strings.TrimSpace(tagName) == "" {
		err = errors.New("slice struct's FieldName can't be empty")
		return
	}
	resultVal := reflect.ValueOf(dataA)
	// use GetIndex
	valueB := reflect.ValueOf(dataB)
	dataBFieldIndex := getStructTagIndex(valueB.Type().Elem(), tagName)
	if dataBFieldIndex < 0 {
		err = errors.New("field `" + tagName + "` not exist in struct")
		return
	}
	var isExist bool = false
	for index := 0; index < valueB.Len(); index++ {
		underlyFieldValueB := reflect.Indirect(valueB.Index(index).Field(dataBFieldIndex))
		if isExist, err = Contains(dataA, underlyFieldValueB.Interface(), tagName); err != nil {
			err = errors.Wrap(err, "getStruct")
			return
		} else if !isExist {
			resultVal = reflect.Append(resultVal, valueB.Index(index))
		}
	}
	return resultVal.Interface(), nil
}

// get index of the field name in  struct data
func getStructTagIndex(typ reflect.Type, tagName string) int {
	for index := 0; index < typ.NumField(); index++ {
		if typ.Field(index).Name == tagName {
			return index
		}
	}
	return -1
}

// interaction
type interaction struct{}

func (t *interaction) getNonStruct(dataA interface{}, dataB interface{}) (result interface{}, err error) {
	valueB := reflect.ValueOf(dataB)
	// new zero value
	resultVal := reflect.MakeSlice(reflect.ValueOf(dataA).Type(), 0, 0)
	valueA := reflect.ValueOf(dataA)
	for index := 0; index < valueB.Len(); index++ {
		var subIndex int = 0
		for subIndex = 0; subIndex < valueA.Len(); subIndex++ {
			aElem := reflect.Indirect(valueA.Index(subIndex))
			bElem := reflect.Indirect(valueB.Index(index))
			if aElem.Interface() == bElem.Interface() {
				resultVal = reflect.Append(resultVal, valueB.Index(index))
				break
			}
		}
	}
	return resultVal.Interface(), nil
}

func (t *interaction) getStruct(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	if strings.TrimSpace(tagName) == "" {
		err = errors.New("slice struct's FieldName can't be empty")
		return
	}
	// new zero value
	resultVal := reflect.MakeSlice(reflect.ValueOf(dataA).Type(), 0, 0)
	// use GetIndex
	valueB := reflect.ValueOf(dataB)
	dataBFieldIndex := getStructTagIndex(valueB.Type().Elem(), tagName)
	if dataBFieldIndex < 0 {
		err = errors.New("field `" + tagName + "` not exist in struct")
		return
	}
	var isExist bool = false
	for index := 0; index < valueB.Len(); index++ {
		underlyFieldValueB := reflect.Indirect(valueB.Index(index).Field(dataBFieldIndex))
		if isExist, err = Contains(dataA, underlyFieldValueB.Interface(), tagName); err != nil {
			err = errors.Wrap(err, "getStruct")
			return
		} else if isExist {
			resultVal = reflect.Append(resultVal, valueB.Index(index))
		}
	}
	return resultVal.Interface(), nil
}

type difference struct{}

func (t *difference) getNonStruct(dataA interface{}, dataB interface{}) (result interface{}, err error) {
	valueB := reflect.ValueOf(dataB)
	// new zero value
	resultVal := reflect.MakeSlice(reflect.ValueOf(dataA).Type(), 0, 0)
	valueA := reflect.ValueOf(dataA)
	for index := 0; index < valueA.Len(); index++ {
		var subIndex int = 0
		for subIndex = 0; subIndex < valueB.Len(); subIndex++ {
			aElem := reflect.Indirect(valueA.Index(index))
			bElem := reflect.Indirect(valueB.Index(subIndex))
			if aElem.Interface() == bElem.Interface() {
				break
			}
		}
		if subIndex == valueB.Len() {
			resultVal = reflect.Append(resultVal, valueA.Index(index))
		}
	}
	return resultVal.Interface(), nil
}

func (t *difference) getStruct(dataA interface{}, dataB interface{}, tagName string) (result interface{}, err error) {
	if strings.TrimSpace(tagName) == "" {
		err = errors.New("slice struct's FieldName can't be empty")
		return
	}
	// new zero value
	resultVal := reflect.MakeSlice(reflect.ValueOf(dataA).Type(), 0, 0)
	// use GetIndex
	valueA := reflect.ValueOf(dataA)
	dataAFieldIndex := getStructTagIndex(valueA.Type().Elem(), tagName)
	if dataAFieldIndex < 0 {
		err = errors.New("field `" + tagName + "` not exist in struct")
		return
	}
	var isExist bool = false
	for index := 0; index < valueA.Len(); index++ {
		underlyFieldValueA := reflect.Indirect(valueA.Index(index).Field(dataAFieldIndex))
		if isExist, err = Contains(dataB, underlyFieldValueA.Interface(), tagName); err != nil {
			err = errors.Wrap(err, "getStruct")
			return
		} else if !isExist {
			resultVal = reflect.Append(resultVal, valueA.Index(index))
		}
	}
	return resultVal.Interface(), nil
}
