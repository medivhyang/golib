package binding

import (
	"reflect"
	"strconv"
	"time"
)

type Binder interface {
	Match(v reflect.Value) bool
	Bind(v reflect.Value, s string) error
}

type BaseBinder struct{}

func (b *BaseBinder) Match(value reflect.Value) bool {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
		return true
	default:
		return false
	}
}

func (b *BaseBinder) Bind(value reflect.Value, s string) error {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.String:
		value.SetString(s)
	case reflect.Int:
		return b.bindInt(value, s, 10, 0)
	case reflect.Int8:
		return b.bindInt(value, s, 10, 8)
	case reflect.Int16:
		return b.bindInt(value, s, 10, 16)
	case reflect.Int32:
		return b.bindInt(value, s, 10, 32)
	case reflect.Int64:
		return b.bindInt(value, s, 10, 64)
	case reflect.Uint:
		return b.bindInt(value, s, 10, 0)
	case reflect.Uint8:
		return b.bindInt(value, s, 10, 8)
	case reflect.Uint16:
		return b.bindUint(value, s, 10, 16)
	case reflect.Uint32:
		return b.bindUint(value, s, 10, 32)
	case reflect.Uint64:
		return b.bindUint(value, s, 10, 64)
	case reflect.Bool:
		return b.bindBool(value, s)
	case reflect.Float32:
		return b.bindFloat(value, s, 32)
	case reflect.Float64:
		return b.bindFloat(value, s, 64)
	}
	return nil
}

func (b *BaseBinder) bindInt(rv reflect.Value, v string, base int, size int) error {
	if v == "" {
		v = "0"
	}
	i, err := strconv.ParseInt(v, base, size)
	if err != nil {
		return err
	}
	rv.SetInt(i)
	return nil
}

func (b *BaseBinder) bindUint(rv reflect.Value, v string, base int, size int) error {
	if v == "" {
		v = "0"
	}
	ui, err := strconv.ParseUint(v, base, size)
	if err != nil {
		return err
	}
	rv.SetUint(ui)
	return nil
}

func (b *BaseBinder) bindFloat(rv reflect.Value, v string, size int) error {
	if v == "" {
		v = "0"
	}
	f, err := strconv.ParseFloat(v, size)
	if err != nil {
		return err
	}
	rv.SetFloat(f)
	return nil
}

func (b *BaseBinder) bindBool(rv reflect.Value, v string) error {
	if v == "" {
		v = "0"
	}
	bl, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	rv.SetBool(bl)
	return nil
}

type TimeBinder struct {
	Layout   string
	Location *time.Location
}

func (b *TimeBinder) Match(rv reflect.Value) bool {
	_, ok := rv.Interface().(time.Time)
	if ok {
		return true
	}
	_, ok2 := rv.Interface().(*time.Time)
	return ok2
}

func (b *TimeBinder) Bind(rv reflect.Value, v string) error {
	layout := time.RFC3339
	if b.Layout != "" {
		layout = time.RFC3339
	}
	location := time.Local
	if b.Location != nil {
		location = b.Location
	}
	t, err := time.ParseInLocation(layout, v, location)
	if err != nil {
		return err
	}
	switch rv.Interface().(type) {
	case time.Time:
		rv.Set(reflect.ValueOf(t))
	case *time.Time:
		rv.Set(reflect.ValueOf(&t))
	}
	return nil
}
