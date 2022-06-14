package main

import (
	"fmt"

	"github.com/medivhyang/golib/string/converting"
)

func main() {
	fmt.Println(converting.Convert[int]("123"))
	fmt.Println(converting.ConvertOrDefault("abc", 456))
}
