package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func BulkInsert(dialect Dialect, models interface{}) TemplateWithError {
	value := reflect.ValueOf(models)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Slice {
		return TemplateWithError{Err: ErrRequireSliceType}
	}
	l := value.Len()
	if l == 0 {
		return TemplateWithError{}
	}
	var matrix [][]ColumnValuePair
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		pairs, err := ParseColumnValuePairs(dialect, elem.Interface())
		if err != nil {
			return TemplateWithError{Err: err}
		}
		if len(pairs) > 0 {
			matrix = append(matrix, pairs)
		}
	}
	if len(matrix) == 0 {
		return TemplateWithError{}
	}
	table := ParseTableName(value.Index(0).Interface())
	columns := make([]string, 0, len(matrix[0]))
	for _, c := range matrix[0] {
		columns = append(columns, c.Column)
	}
	holders := make([]string, 0, len(columns))
	for i := 0; i < len(matrix); i++ {
		holders = append(holders, fmt.Sprintf("(%s)", strings.Join(repeatString("?", len(columns)), ",")))
	}
	values := make([]interface{}, 0, len(columns)*len(matrix))
	for _, cc := range matrix {
		for _, c := range cc {
			values = append(values, c.Value)
		}
	}
	t := NewTemplate(fmt.Sprintf("insert into %s(%s) values %s",
		table,
		strings.Join(columns, ", "),
		strings.Join(holders, ", "),
	), values...)
	return TemplateWithError{Template: t}
}
