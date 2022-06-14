package orm

import (
	"fmt"
	"reflect"
)

func ExampleBulkInsert() {
	type User struct {
		Name string `orm:"name"`
		Age  int    `orm:"age"`
	}
	t := BulkInsert(&TestDialect{}, []User{
		{"Medivh", 18},
		{"Jason", 22},
		{"Mike", 30},
	})
	fmt.Println(t)
	// output:
	// "insert into User(name, age) values (?,?), (?,?), (?,?)": []interface {}{"Medivh", 18, "Jason", 22, "Mike", 30}
}

type TestDialect struct{}

func (d *TestDialect) MappingType(rt reflect.Type) string {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Bool:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	default:
		switch rt.String() {
		case "time.Time":
			return "datetime"
		}
		return ""
	}
}

func (d *TestDialect) Quote(s string) string {
	return fmt.Sprintf("'%s'", s)
}
