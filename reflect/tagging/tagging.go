package tagging

import (
	"errors"
	"reflect"
	"strings"
)

var (
	errorPrefix = "tagging: "

	ErrRequireStructType = errors.New(errorPrefix + "require struct type")
)

func ParseStructTag(i interface{}, key string) map[string]string {
	result := map[string]string{}
	m := ParseStructTags(i, key)
	for k, vm := range m {
		if v, ok := vm[key]; ok {
			result[k] = v
		}
	}
	return result
}

func ParseStructTags(i interface{}, keys ...string) map[string]map[string]string {
	if len(keys) == 0 {
		return map[string]map[string]string{}
	}
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(ErrRequireStructType)
	}
	result := map[string]map[string]string{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		rf := t.Field(i)
		result[rf.Name] = map[string]string{}
		for _, k := range keys {
			if v, ok := rf.Tag.Lookup(k); ok {
				result[rf.Name][k] = v
			}
		}
	}
	return result
}

func ParseStructChildTags(i interface{}, parentKey string, itemSep string, kvSep string, keys ...string) (map[string]map[string]string, error) {
	rv := reflect.ValueOf(i)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, ErrRequireStructType
	}
	result := map[string]map[string]string{}
	rt := reflect.TypeOf(i)
	for i := 0; i < rt.NumField(); i++ {
		var (
			field = rt.Field(i)
			tag   = field.Tag.Get(parentKey)
			items []string
			kv    []string
			k, v  string
		)
		result[field.Name] = map[string]string{}
		items = strings.Split(tag, itemSep)
		for _, item := range items {
			kv = strings.Split(item, kvSep)
			if len(kv) >= 1 {
				k = kv[0]
			}
			if len(kv) >= 2 {
				v = kv[1]
			}
			if len(keys) == 0 || containString(keys, k) {
				result[field.Name][k] = v
			}
		}
	}
	return result, nil
}

func containString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}