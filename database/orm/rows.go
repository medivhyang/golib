package orm

import (
	"database/sql"
	"reflect"
)

type Rows struct {
	dialect Dialect
	raw     *sql.Rows
}

func NewRows(dialect Dialect, rows *sql.Rows) *Rows {
	return &Rows{dialect: dialect, raw: rows}
}

func (r *Rows) Scan(callback func(scan func(...interface{}) error, abort func()) error) error {
	if callback == nil {
		return nil
	}
	abort := false
	abortFunc := func() { abort = true }
	for r.raw.Next() {
		if err := callback(r.raw.Scan, abortFunc); err != nil {
			return err
		}
		if abort {
			break
		}
	}
	return r.raw.Close()
}

func (r *Rows) Bind(dst interface{}) error {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Ptr {
		return ErrRequirePointerType
	}
	value = unrefValueAndInit(value)
	switch value.Interface().(type) {
	case map[string]interface{}:
		if !value.CanSet() {
			return ErrCannotSetValue
		}
		m, err := r.Map()
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(m))
	default:
		switch value.Kind() {
		case reflect.Struct:
			if err := r.Struct(value); err != nil {
				return err
			}
		case reflect.Slice:
			switch value.Type().Elem().Kind() {
			case reflect.Map:
				ss, err := r.MapSlice()
				if err != nil {
					return err
				}
				value.Set(reflect.ValueOf(ss))
			case reflect.Struct:
				if err := r.StructSlice(value); err != nil {
					return err
				}
			default:
				if err := r.ScalarSlice(value); err != nil {
					return err
				}
			}
		default:
			if err := r.Scalar(value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Rows) Scalar(value interface{}) error {
	if !r.raw.Next() {
		return sql.ErrNoRows
	}
	if err := r.raw.Scan(value); err != nil {
		return err
	}
	return r.raw.Close()
}

func (r *Rows) ScalarSlice(slice interface{}) error {
	reflectValue := reflect.ValueOf(slice)
	if reflectValue.Kind() != reflect.Ptr {
		return ErrRequirePointerType
	}
	reflectValue = reflectValue.Elem()
	if reflectValue.Kind() != reflect.Slice {
		return errorf("scan rows to values: require slice type")
	}

	elemType := reflectValue.Type().Elem()

	for r.raw.Next() {
		item := reflect.New(elemType)
		if err := r.raw.Scan(item.Interface()); err != nil {
			return err
		}
		reflectValue.Set(reflect.Append(reflectValue, item.Elem()))
	}

	return r.raw.Close()
}

func (r *Rows) Map() (map[string]interface{}, error) {
	columns, err := r.raw.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	if !r.raw.Next() {
		return nil, sql.ErrNoRows
	}
	if err := r.raw.Scan(values...); err != nil {
		return nil, err
	}

	item := make(map[string]interface{})
	for i, v := range values {
		item[columns[i]] = *v.(*interface{})
	}
	if err := r.raw.Close(); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Rows) MapSlice() ([]map[string]interface{}, error) {
	columns, err := r.raw.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	var items []map[string]interface{}
	for r.raw.Next() {
		if err := r.raw.Scan(values...); err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		for i, v := range values {
			item[columns[i]] = *v.(*interface{})
		}
		items = append(items, item)
	}

	if err := r.raw.Close(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Rows) Struct(i interface{}) error {
	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr {
		return ErrRequirePointerType
	}
	columns, err := r.raw.Columns()
	if err != nil {
		return err
	}
	columnValuePairs, err := ParseColumnValuePairs(r.dialect, i)
	if err != nil {
		return err
	}
	columnValuePairsMap := make(map[string]interface{}, len(columnValuePairs))
	for _, pair := range columnValuePairs {
		columnValuePairsMap[pair.Column] = pair.Value
	}
	vv := make([]interface{}, 0, len(columns))
	for _, c := range columns {
		if f, ok := columnValuePairsMap[c]; ok {
			vv = append(vv, f)
		} else {
			var tmp interface{}
			vv = append(vv, &tmp)
		}
	}
	if !r.raw.Next() {
		return sql.ErrNoRows
	}
	if err := r.raw.Scan(vv...); err != nil {
		return err
	}
	if err := r.raw.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Rows) StructSlice(i interface{}) error {
	reflectValue := reflect.ValueOf(i)
	if reflectValue.Kind() != reflect.Ptr {
		return ErrRequirePointerType
	}
	unrefReflectValue := unrefValueAndInit(reflectValue)
	if unrefReflectValue.Kind() != reflect.Slice {
		return ErrRequireSliceType
	}
	columns, err := r.raw.Columns()
	if err != nil {
		return err
	}
	values := make([]interface{}, 0, len(columns))
	for r.raw.Next() {
		item := reflect.New(unrefReflectValue.Type().Elem())
		columnValuePairs, err := ParseColumnValuePairs(r.dialect, item)
		if err != nil {
			return err
		}
		columnValuePairsMap := make(map[string]interface{}, len(columnValuePairs))
		for _, pair := range columnValuePairs {
			columnValuePairsMap[pair.Column] = pair.Value
		}
		for _, c := range columns {
			if v, ok := columnValuePairsMap[c]; ok {
				values = append(values, v)
			} else {
				var temp interface{}
				values = append(values, &temp)
			}
		}
		if err := r.raw.Scan(values...); err != nil {
			return err
		}
		reflectValue.Set(reflect.Append(reflectValue, item))
	}
	if err := r.raw.Close(); err != nil {
		return err
	}
	return nil
}
