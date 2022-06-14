package binding

import (
	"errors"
	"reflect"
)

var (
	ErrUnsupportedType = errors.New("binding: unsupported type")
	ErrTooMuchValues   = errors.New("binding: too much values")
)

var DefaultBinders = []Binder{&BaseBinder{}, &TimeBinder{}}

func Bind(src string, dst interface{}, binders ...Binder) error {
	if len(binders) == 0 {
		binders = DefaultBinders
	}
	return bind(src, reflect.ValueOf(dst), binders...)
}

func bind(src string, dst reflect.Value, binders ...Binder) error {
	for _, b := range binders {
		if b.Match(dst) {
			return b.Bind(dst, src)
		}
	}
	return ErrUnsupportedType
}

func BindList(src []string, dst interface{}, binders ...Binder) error {
	if len(binders) == 0 {
		binders = DefaultBinders
	}
	return bindList(src, reflect.ValueOf(dst), binders...)
}

func bindList(src []string, dst reflect.Value, binders ...Binder) error {
	drv := unrefValueAndInit(dst)
	switch unrefType(dst.Type()).Kind() {
	case reflect.Slice:
		drv.Set(reflect.MakeSlice(drv.Type(), len(src), len(src)))
		for i, v := range src {
			if err := bind(v, drv.Index(i), binders...); err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		if len(src) >= drv.Len() {
			return ErrTooMuchValues
		}
		for i, v := range src {
			return bind(v, drv.Field(i), binders...)
		}
		return nil
	default:
		return ErrUnsupportedType
	}
}

func BindStruct(src map[string][]string, dst interface{}, binders ...Binder) error {
	if len(binders) == 0 {
		binders = DefaultBinders
	}
	return bindStruct(src, reflect.ValueOf(dst), binders...)
}

func bindStruct(src map[string][]string, dst reflect.Value, binders ...Binder) error {
	return bindStructFunc(func(s string) []string {
		return src[s]
	}, dst, binders...)
}

func BindStructFunc(src func(string) []string, dst interface{}, binders ...Binder) error {
	if len(binders) == 0 {
		binders = DefaultBinders
	}
	return bindStructFunc(src, reflect.ValueOf(dst), binders...)
}

func bindStructFunc(src func(string) []string, dst reflect.Value, binders ...Binder) error {
	if src == nil {
		return nil
	}
	if len(binders) == 0 {
		binders = DefaultBinders
	}
	for dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	switch dst.Kind() {
	case reflect.Struct:
		for i := 0; i < dst.NumField(); i++ {
			if isUnexportedStructField(dst.Type().Field(i)) {
				continue
			}
			var (
				fv   = dst.Field(i)
				name = dst.Type().Field(i).Name
				vv   = src(name)
			)
			if len(vv) == 0 {
				continue
			}
			switch unrefType(fv.Type()).Kind() {
			case reflect.Array, reflect.Slice:
				if err := bindList(vv, fv, binders...); err != nil {
					return err
				}
			default:
				v := ""
				if len(vv) > 0 {
					v = vv[0]
				}
				if err := bind(v, fv, binders...); err != nil {
					return err
				}
			}
		}
	default:
		return ErrUnsupportedType
	}
	return nil
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

func unrefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func isUnexportedStructField(sf reflect.StructField) bool {
	return sf.PkgPath != "" && !sf.Anonymous
}