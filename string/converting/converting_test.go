package converting

import (
	"fmt"
	"time"
)

func ExampleConvert() {
	fmt.Println(Convert[int64]("123"))
	fmt.Println(Convert[rune]("123"))
	fmt.Println(Convert[float64]("0.123"))
	fmt.Println(Convert[time.Time]("2022-06-14T16:11:00Z"))
	fmt.Println(Convert[time.Time]("2022-06-14 16:11:00", "2006-01-02 15:04:05", time.Local))
	fmt.Println(Convert[time.Duration]("1m30s"))

	// Output:
	// 123 <nil>
	// 123 <nil>
	// 0.123 <nil>
	// 2022-06-14 16:11:00 +0000 UTC <nil>
	// 2022-06-14 16:11:00 +0800 CST <nil>
	// 1m30s <nil>
}

func ExampleConvertOrDefault() {
	fmt.Println(ConvertOrDefault("not an int", 123))
	fmt.Println(ConvertOrDefault("not a time", time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// Output:
	// 123
	// 2020-01-01 00:00:00 +0000 UTC
}

func ExampleBatchConvert() {
	fmt.Println(BatchConvert[int]([]string{"1", "2", "3"}))
	fmt.Println(BatchConvert[byte]([]string{"1", "2", "3"}))
	fmt.Println(BatchConvert[time.Time]([]string{"2022-06-14 16:11:00", "2022-06-14 16:12:00", "2022-06-14 16:13:00"}))

	// Output:
	// [1 2 3] <nil>
	// [1 2 3] <nil>
	// [2022-06-14 16:11:00 +0000 UTC 2022-06-14 16:12:00 +0000 UTC 2022-06-14 16:13:00 +0000 UTC] <nil>
}

func ExampleBatchConvertOrDefault() {
	fmt.Println(BatchConvertOrDefault([]string{"not a int"}, []int{1, 2, 3}))

	// Output:
	// [1 2 3]
}
