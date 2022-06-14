package binding

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestBindStruct(t *testing.T) {
	type args struct {
		input   map[string][]string
		output  reflect.Value
		binders []Binder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				output: reflect.ValueOf(&struct {
					Name    string
					Age     int
					Pets    *[]string
					Created *time.Time
					hidden  string
				}{}),
				input: map[string][]string{
					"Name":    {"Medivh", "Mike"},
					"Age":     {"28"},
					"Pets":    {"Foo", "Bar", "Baz"},
					"Created": {"2020-02-22T22:22:22Z"},
					"hidden":  {"secret"},
				},
				binders: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bindStruct(tt.args.input, tt.args.output, tt.args.binders...); (err != nil) != tt.wantErr {
				t.Errorf("BindStructDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%+v", tt.args.output.Interface())
		})
	}
}

func TestBindList(t *testing.T) {
	type args struct {
		output  reflect.Value
		input   []string
		binders []Binder
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				output:  reflect.ValueOf(&[]int{}),
				input:   []string{"1", "2", "3"},
				binders: nil,
			},
			wantErr: false,
		},
		{
			name: "parse fail",
			args: args{
				output: reflect.ValueOf(&[]int{}),
				input:  []string{"1", "2", "3", "a"}, binders: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bindList(tt.args.input, tt.args.output, tt.args.binders...); (err != nil) != tt.wantErr {
				t.Errorf("BindListDefault() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%+v", tt.args.output.Interface())
		})
	}
}

func ExampleBind() {
	var output int
	if err := Bind("123", &output); err != nil {
		panic(err)
	}
	fmt.Println(output)
	// output:
	// 123
}

func ExampleBindList() {
	var output []int
	if err := BindList([]string{"1", "2", "3"}, &output); err != nil {
		panic(err)
	}
	fmt.Println(output)
	// output:
	// [1 2 3]
}
