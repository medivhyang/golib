package main

import (
	"fmt"

	"github.com/medivhyang/golib/container/slices"
)

func main() {
	fmt.Println(slices.Map([]int{1, 2, 3, 4, 5}, func(a int) int {
		return a * 2
	}))
	fmt.Println(slices.Contains([]int{1, 2, 3}, 1, 2, 3, 4, 5))
}
