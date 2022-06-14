package orm

import (
	"fmt"
)

func ExampleCreateTablesIfNotExists() {
	type User struct {
		Name string `orm:"name varchar(512) primary key"`
		Age  int    `orm:"age"`
	}

	t := CreateTablesIfNotExists(&TestDialect{}, User{})
	fmt.Println(t)

	// output:
	// "create table if not exists User ('name' varchar(512) primary key, 'age' integer);"
}

func ExampleDropTablesIfExists() {
	type User struct {
		Name string `orm:"name varchar(512) primary key"`
		Age  int    `orm:"age"`
	}

	t := DropTablesIfExists(&TestDialect{}, User{})
	fmt.Println(t)

	// output:
	// "drop table if exists User;"
}
