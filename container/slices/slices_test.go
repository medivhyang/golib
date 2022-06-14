package slices

import "fmt"

func ExampleFilter() {
	var list = []int{1, 2, 3, 4, 5}
	var result = Filter(list, func(a int) bool {
		return a > 2
	})
	fmt.Println(result)

	// Output:
	// [3 4 5]
}

func ExampleMap() {
	var list = []int{1, 2, 3, 4, 5}
	var result = Map(list, func(a int) int {
		return a * 2
	})
	fmt.Println(result)

	// Output:
	// [2 4 6 8 10]
}

func ExampleReduce() {
	var list = []int{1, 2, 3, 4, 5}
	var sum = Reduce(list, func(a, b int) int {
		return a + b
	})
	fmt.Println(sum)

	// Output:
	// 15
}
