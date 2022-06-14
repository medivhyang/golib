package orm

import (
	"reflect"
	"sync"
)

const DefaultDialectName = ""

type Dialect interface {
	MappingType(rt reflect.Type) string
	Quote(s string) string
}

var dialects sync.Map

func RegisterDialect(name string, dialect Dialect) {
	dialects.Store(name, dialect)
}

func GetDialect(name string) Dialect {
	v, ok := dialects.Load(name)
	if !ok {
		panic(ErrNotFoundDialect)
	}
	d, ok := v.(Dialect)
	if !ok {
		panic(ErrInvalidDialect)
	}
	return d
}

func LookupDialect(name string) (Dialect, bool) {
	v, ok := dialects.Load(name)
	if !ok {
		return nil, false
	}
	d, ok := v.(Dialect)
	if !ok {
		return nil, false
	}
	return d, true
}

func RegisterDefaultDialect(name string, dialect Dialect) {
	RegisterDialect(name, dialect)
	SetDefaultDialect(name)
}

func GetDefaultDialect() Dialect {
	return GetDialect(DefaultDialectName)
}

func SetDefaultDialect(name string) bool {
	d, ok := LookupDialect(name)
	if ok {
		dialects.Store(DefaultDialectName, d)
		return true
	}
	return false
}
