package diff

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type DiffValue struct {
	Name   string
	Before interface{}
	After  interface{}
}

func (dv DiffValue) String() string {
	return dv.Name + ": " + Strval(dv.Before) + " -> " + Strval(dv.After)
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func DiffStructStr[T any](old, new T) []string {
	diffValues := DiffStruct(old, new)
	diffStr := make([]string, 0, len(diffValues))
	for _, v := range diffValues {
		diffStr = append(diffStr, v.String())
	}
	return diffStr
}

func DiffStruct[T any](old, new T) []DiffValue {

	types := reflect.TypeOf(old)
	valOld := reflect.ValueOf(old)
	valNew := reflect.ValueOf(new)
	if types.Kind() == reflect.Ptr {
		types = reflect.TypeOf(old).Elem()
		valOld = valOld.Elem()
		valNew = valNew.Elem()
	}

	results := make([]DiffValue, 0)

	for i := 0; i < types.NumField(); i++ {
		tag, hasTag := readTag(types.Field(i))
		if !hasTag {
			continue
		}
		fieldName := types.Field(i).Name
		if valOld.Field(i).Interface() == valNew.FieldByName(fieldName).Interface() {
			continue
		}
		if valOld.Field(i).IsZero() && valNew.FieldByName(fieldName).IsZero() {
			continue
		}

		results = append(results, DiffValue{
			Name:   tag,
			Before: valOld.Field(i).Interface(),
			After:  valNew.FieldByName(fieldName).Interface(),
		})
	}
	return results
}

func readTag(f reflect.StructField) (string, bool) {
	val, ok := f.Tag.Lookup("diff")

	// no tag, skip this field
	if val == "" || !ok {
		return "", false
	}

	return val, true
}
