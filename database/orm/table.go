package orm

import (
	"reflect"
	"strings"
)

type Table struct {
	Name    string
	Columns []Column
}

func (t *Table) ColumnNames() ([]string, error) {
	var names []string
	for _, c := range t.Columns {
		names = append(names, c.Name)
	}
	return names, nil
}

type Column struct {
	Name   string
	Type   string
	Suffix string
}

type ColumnValuePair struct {
	Column string
	Value  interface{}
}

func ParseTables(dialect Dialect, models ...interface{}) ([]*Table, error) {
	tt := make([]*Table, 0, len(models))
	for _, model := range models {
		t, err := ParseTable(dialect, model)
		if err != nil {
			return nil, err
		}
		if t != nil {
			tt = append(tt, t)
		}
	}
	return tt, nil
}

func ParseTable(dialect Dialect, model interface{}) (*Table, error) {
	if dialect == nil {
		dialect = GetDefaultDialect()
	}
	value := reflect.ValueOf(model)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, ErrRequireStructType
	}
	t := Table{Name: ParseTableName(model), Columns: nil}
	for i := 0; i < value.NumField(); i++ {
		c, err := ParseColumn(dialect, value.Type().Field(i))
		if err != nil {
			return nil, err
		}
		t.Columns = append(t.Columns, c)
	}
	return &t, nil
}

func ParseTableName(model interface{}) string {
	obj, ok := model.(interface{ Table() string })
	if !ok {
		return reflect.TypeOf(model).Name()
	}
	return obj.Table()
}

func ParseColumn(dialect Dialect, sf reflect.StructField) (Column, error) {
	if dialect == nil {
		dialect = GetDefaultDialect()
	}
	items := strings.Split(sf.Tag.Get(TagKey), " ")

	var finalName string
	if len(items) >= 1 {
		finalName = items[0]
	}
	if finalName == "" {
		finalName = sf.Name
	}
	finalName = toCase(caseSnake, finalName)

	var finalType string
	if len(items) >= 2 {
		finalType = items[1]
	}
	if finalType == "" {
		finalType = dialect.MappingType(sf.Type)
	}

	var suffix string
	if len(items) > 2 {
		suffix = strings.Join(items[2:], " ")
	}

	return Column{
		Name:   finalName,
		Type:   finalType,
		Suffix: suffix,
	}, nil
}

func ParseColumnValuePairs(dialect Dialect, model interface{}) ([]ColumnValuePair, error) {
	value := reflect.ValueOf(model)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, ErrRequireStructType
	}

	var pairs []ColumnValuePair
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)
		// check struct field is unexported
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}
		c, err := ParseColumn(dialect, sf)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, ColumnValuePair{
			Column: c.Name,
			Value:  value.Field(i).Interface(),
		})
	}

	return pairs, nil
}
