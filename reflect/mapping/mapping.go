package mapping

import (
	"errors"
	"reflect"
)

var (
	ErrRequireStructType  = errors.New("mapping: require struct type")
	ErrRequirePointerType = errors.New("mapping: require pointer type")
)

func MapStructToStruct(src interface{}, dst interface{}, nameFunc func(string) string) error {
	srcValue := reflect.ValueOf(src)
	if unrefType(srcValue.Type()).Kind() != reflect.Struct {
		return ErrRequireStructType
	}
	dstValue := reflect.ValueOf(dst)
	if dstValue.Type().Kind() != reflect.Ptr {
		return ErrRequirePointerType
	}
	if unrefType(dstValue.Type()).Kind() != reflect.Struct {
		return ErrRequireStructType
	}
	result := MapStructToMap(src)
	copyMap := make(map[string]interface{}, len(result))
	for k, v := range result {
		copyMap[k] = v
	}
	for name, value := range copyMap {
		newName := nameFunc(name)
		if newName != name {
			result[newName] = value
			delete(result, name)
		}
	}
	return MapMapToStruct(result, dst, nil)
}

func MapMapToStruct(src map[string]interface{}, dst interface{}, nameFunc func(string) string) error {
	dstValue := reflect.ValueOf(dst)
	dstType := dstValue.Type()
	if dstType.Kind() != reflect.Ptr {
		panic(ErrRequirePointerType)
	}
	if unrefType(dstType).Kind() != reflect.Struct {
		panic(ErrRequireStructType)
	}
	dstValue = unrefValueAndInit(dstValue)
	dstMap := map[string]reflect.Value{}
	for i := 0; i < dstValue.NumField(); i++ {
		dstMap[dstValue.Type().Field(i).Name] = dstValue.Field(i)
	}
	for name, srcFieldValue := range src {
		if nameFunc != nil {
			name = nameFunc(name)
		}
		dstFieldValue, ok := dstMap[name]
		if !ok {
			continue
		}
		dstFieldType := dstFieldValue.Type()
		srcReflectFieldValue := reflect.ValueOf(srcFieldValue)
		srcReflectFieldType := srcReflectFieldValue.Type()
		if unrefType(srcReflectFieldType) != unrefType(dstFieldValue.Type()) {
			continue
		}
		switch unrefType(srcReflectFieldType).Kind() {
		case reflect.Slice:
			l := unrefValue(srcReflectFieldValue).Len()
			unrefValueAndInit(dstFieldValue).Set(reflect.MakeSlice(unrefType(dstFieldType), l, l))
			fallthrough
		case reflect.Array:
			reflect.Copy(unrefValueAndInit(dstFieldValue), unrefValueAndInit(srcReflectFieldValue))
		default:
			dstFieldValue.Set(reflect.ValueOf(srcFieldValue))
		}
	}
	return nil
}

func MapStructToMap(src interface{}) map[string]interface{} {
	value := unrefValue(reflect.ValueOf(src))
	if value.Kind() != reflect.Struct {
		panic(ErrRequireStructType)
	}
	dst := map[string]interface{}{}
	for i := 0; i < value.NumField(); i++ {
		dst[value.Type().Field(i).Name] = value.Field(i).Interface()
	}
	return dst
}

func unrefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func unrefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func unrefValueAndInit(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}
