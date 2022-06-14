package copying

import (
	"fmt"
	"net/http"
	"time"
)

func ExampleDeepCopy() {
	src := &http.Client{Timeout: 100 * time.Second}
	dst := &http.Client{}
	if err := DeepCopy(dst, src); err != nil {
		fmt.Println(err)
	}
	fmt.Println(src)
	fmt.Println(dst)

	// Output:
	// &{<nil> <nil> <nil> 1m40s}
	// &{<nil> <nil> <nil> 1m40s}
}
