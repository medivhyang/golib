package orm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Template struct {
	Format string
	Values []interface{}
}

func NewTemplate(format string, values ...interface{}) Template {
	t := Template{Format: format, Values: append([]interface{}{}, values...)}
	return t
}

func (t Template) Merge(others ...Template) Template {
	newT := NewTemplate(t.Format, t.Values...)
	for _, o := range others {
		newT.Format += o.Format
		newT.Values = append(o.Values, o.Values...)
	}
	return newT
}

func (t Template) Appendf(format string, values ...interface{}) Template {
	return NewTemplate(t.Format+format, append(append([]interface{}{}, t.Values...), values...)...)
}

func (t Template) AppendValues(values ...interface{}) Template {
	if len(t.Values) > 0 {
		t.Values = append(append([]interface{}{}, t.Values...), values...)
	}
	return t
}

func (t Template) Wrap(left string, right string) Template {
	t.Format = left + t.Format + right
	return t
}

func (t Template) Bracket() Template {
	return t.Wrap("(", ")")
}

func (t Template) IsEmpty() bool {
	return t.Format == "" && len(t.Values) == 0
}

func (t Template) Exec(ctx context.Context, db DBTX) (sql.Result, error) {
	return db.Exec(ctx, t)
}

func (t Template) Query(ctx context.Context, db DBTX, value interface{}) error {
	return db.Query(ctx, t, value)
}

func (t Template) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%q", t.Format))
	if len(t.Values) > 0 {
		b.WriteString(": ")
		b.WriteString(fmt.Sprintf("%#v", t.Values))
	}
	return b.String()
}

type TemplateWithError struct {
	Template
	Err error
}

func (t TemplateWithError) Exec(ctx context.Context, db DBTX) (sql.Result, error) {
	if t.Err != nil {
		return nil, t.Err
	}
	return db.Exec(ctx, t.Template)
}

func (t TemplateWithError) Query(ctx context.Context, db DBTX, i interface{}) error {
	if t.Err != nil {
		return t.Err
	}
	return db.Query(ctx, t.Template, i)
}

type Templates []Template

func (tt Templates) Append(others ...Template) Templates {
	return append(tt, others...)
}

func (tt Templates) Appendf(format string, values ...interface{}) Templates {
	return tt.Append(NewTemplate(format, values...))
}

func (tt Templates) Join(sep string, right, left string) Template {
	if len(tt) == 0 {
		return Template{}
	}
	ff := make([]string, 0, len(tt))
	vv := make([]interface{}, 0, len(tt))
	for _, c := range tt {
		ff = append(ff, c.Format)
		vv = append(vv, c.Values...)
	}
	return NewTemplate(right+strings.Join(ff, sep)+left, vv...)
}

func (tt Templates) Formats() []string {
	r := make([]string, 0, len(tt))
	for _, t := range tt {
		r = append(r, t.Format)
	}
	return r
}

func (tt Templates) Values() []interface{} {
	var r []interface{}
	for _, t := range tt {
		r = append(r, t.Values...)
	}
	return r
}
