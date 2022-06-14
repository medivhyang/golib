package config

import "fmt"

func ExampleLoad() {
	src := `{"name":"Medivh", "age":99}`
	dst := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	if err := LoadString(DefaultJSONSource, src, dst); err != nil {
		panic(err)
	}
	fmt.Println(dst.Name, dst.Age)

	// output:
	// Medivh 99
}
