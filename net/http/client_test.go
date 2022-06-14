package http

import (
	"fmt"
)

func ExampleGet() {
	fmt.Println(Get("https://www.baidu.com").StatusCode())

	// Output: 200
}

func ExampleTemplate_New() {
	t := &Template{
		Prefix: "https://example.com/",
		Queries: map[string]string{
			"foo": "bar",
		},
	}
	fmt.Println(t.FullURL())
	fmt.Println(t.New().Get("/user").Query("id", "1").Queries(map[string]string{"name": "medivh"}).Build().FullURL())

	// Output:
	// https://example.com/?foo=bar
	// https://example.com/user?foo=bar&id=1&name=medivh
}
