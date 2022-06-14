package converting

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var defaultLayouts = []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05"}

type convertable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~complex64 | ~complex128 |
		~string | ~bool | time.Time
}

// ConvertOrDefault convert string to some type, if error, return default value
func ConvertOrDefault[T convertable](s string, v T, options ...any) T {
	result, err := Convert[T](s, options...)
	if err != nil {
		return v
	}
	return result
}

// Convert string to some type
func Convert[T convertable](s string, options ...any) (T, error) {
	var result T
	v := reflect.ValueOf(&result).Elem()
	switch v.Type().String() {
	case "time.Time":
		find := false
		switch len(options) {
		case 0:
			for _, l := range defaultLayouts {
				t, err := time.Parse(l, s)
				if err == nil {
					find = true
					v.Set(reflect.ValueOf(t))
				}
			}
			if !find {
				return result, errors.New("time parse error")
			}
		case 1:
			t, err := time.Parse(options[0].(string), s)
			if err != nil {
				return result, err
			}
			v.Set(reflect.ValueOf(t))
		case 2:
			t, err := time.ParseInLocation(options[0].(string), s, options[1].(*time.Location))
			if err != nil {
				return result, err
			}
			v.Set(reflect.ValueOf(t))
		default:
			return result, errors.New("time parse error")
		}
	case "time.Duration":
		d, err := time.ParseDuration(s)
		if err != nil {
			return result, err
		}
		v.Set(reflect.ValueOf(d))
	default:
		switch v.Type().Kind() {
		case reflect.String:
			v.SetString(s)
		case reflect.Int:

			i, err := strconv.Atoi(s)
			if err != nil {
				return result, err
			}
			v.SetInt(int64(i))
		case reflect.Int8:
			i, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return result, err
			}
			v.SetInt(int64(i))
		case reflect.Int16:
			i, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return result, err
			}
			v.SetInt(int64(i))
		case reflect.Int32:
			i, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return result, err
			}
			v.SetInt(int64(i))
		case reflect.Int64:
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return result, err
			}
			v.SetInt(int64(i))
		case reflect.Uint:
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return result, err
			}
			v.SetUint(uint64(i))
		case reflect.Uint8:
			i, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return result, err
			}
			v.SetUint(uint64(i))
		case reflect.Uint16:
			i, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return result, err
			}
			v.SetUint(uint64(i))
		case reflect.Uint32:
			i, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return result, err
			}
			v.SetUint(uint64(i))
		case reflect.Uint64:
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return result, err
			}
			v.SetUint(uint64(i))
		case reflect.Float32:
			i, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return result, err
			}
			v.SetFloat(float64(i))
		case reflect.Float64:
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return result, err
			}
			v.SetFloat(float64(i))
		case reflect.Bool:
			i, err := strconv.ParseBool(s)
			if err != nil {
				return result, err
			}
			v.SetBool(i)
		default:
			return result, errors.New("unsupported type")
		}
	}
	return result, nil
}

// BatchConvertOrDefault batch convert string to some type, if error, return default value
func BatchConvertOrDefault[T convertable](ss []string, v []T, options ...any) []T {
	result, err := BatchConvert[T](ss, options...)
	if err != nil {
		return v
	}
	return result
}

// BatchConvert batch convert string to some type
func BatchConvert[T convertable](ss []string, options ...any) ([]T, error) {
	result := make([]T, len(ss))
	for i, v := range ss {
		v, err := Convert[T](v, options...)
		if err != nil {
			return result, err
		}
		result[i] = v
	}
	return result, nil
}
