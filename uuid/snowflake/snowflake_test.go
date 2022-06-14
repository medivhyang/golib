package snowflake

import (
	"fmt"
)

func ExampleGenerate() {
	fmt.Println(len(Generate().Base36()))

	// output:
	// 12
}
