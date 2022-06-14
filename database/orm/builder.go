package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Action string

const (
	actionSelect Action = "select"
	actionInsert Action = "insert"
	actionUpdate Action = "update"
	actionDelete Action = "delete"
)

type Builder struct {
	dialect  Dialect
	action   Action
	table    Template
	columns  Templates
	joins    Templates
	where    Condition
	orderBy  []string
	paging   Template
	groupBy  []string
	having   Condition
	distinct bool
	err      error
}

func New(dialect ...Dialect) *Builder {
	b := new(Builder)
	if len(dialect) > 0 && dialect[0] != nil {
		b.Dialect(dialect[0])
	} else {
		b.Dialect(GetDefaultDialect())
	}
	return b
}

func (b *Builder) Build() TemplateWithError {
	if b.err != nil {
		return TemplateWithError{Err: b.err}
	}
	t := Template{}
	switch b.action {
	case actionSelect:
		columns := make([]string, 0, len(b.columns.Formats()))
		for _, c := range b.columns.Formats() {
			columns = append(columns, b.dialect.Quote(c))
		}
		if len(columns) == 0 {
			columns = []string{"*"}
		}
		if b.distinct {
			t = t.Appendf(fmt.Sprintf("select distinct %s from %s",
				strings.Join(columns, ","),
				b.dialect.Quote(b.table.Format),
			), append(append([]interface{}{}, b.columns.Values()...), b.table.Values...))
		} else {
			t = t.Appendf(fmt.Sprintf("select %s from %s",
				strings.Join(columns, ","),
				b.dialect.Quote(b.table.Format),
			), b.columns.Values()...)
		}
		if len(b.joins) > 0 {
			t = t.Merge(b.joins.Join(" ", "", ""))
		}
		if b.where.IsNotEmpty() {
			t = t.Appendf(" where ").Merge(b.where.And())
		}
		if len(b.groupBy) > 0 {
			t = t.Appendf(fmt.Sprintf(" group by %s", strings.Join(b.groupBy, ",")))
		}
		if b.having.IsNotEmpty() {
			t = t.Appendf(" having ").Merge(b.having.And())
		}
		if len(b.orderBy) > 0 {
			t = t.Appendf(fmt.Sprintf(" order by %s", strings.Join(b.orderBy, ", ")))
		}
		if len(b.paging.Format) > 0 {
			t = t.Appendf(" ").Merge(b.paging)
		}
	case actionInsert:
		t = t.Appendf(fmt.Sprintf("insert into %s(%s) values(%s)",
			b.table,
			strings.Join(b.columns.Formats(), ", "),
			strings.Join(repeatString("?", len(b.columns)), ","),
		), b.columns.Values()...)
		if b.where.IsNotEmpty() {
			t = t.Appendf(" where ").Merge(b.where.And())
		}
	case actionUpdate:
		pairs := make([]string, 0, len(b.columns))
		for _, c := range b.columns {
			pairs = append(pairs, fmt.Sprintf("%s = ?", b.dialect.Quote(c.Format)))
		}
		t.Appendf(fmt.Sprintf("update %s set %s",
			b.table,
			strings.Join(pairs, ","),
		), b.columns.Values()...)
		if b.where.IsNotEmpty() {
			t = t.Appendf(" where ").Merge(b.where.And())
		}
	case actionDelete:
		t = t.Appendf(fmt.Sprintf("delete from %s", b.table))
		if b.where.IsNotEmpty() {
			t = t.Appendf(" where ").Merge(b.where.And())
		}
	}
	return TemplateWithError{Template: t}
}

func (b *Builder) Query(ctx context.Context, db DBTX, i interface{}) error {
	if b.err != nil {
		return b.err
	}
	return b.Build().Query(ctx, db, i)
}

func (b *Builder) Exec(ctx context.Context, db DBTX) (sql.Result, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.Build().Exec(ctx, db)
}

func (b *Builder) Dialect(d Dialect) *Builder {
	if b.err != nil {
		return b
	}
	b.dialect = d
	return b
}

func (b *Builder) Select(table string, columns ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionSelect
	b.table = NewTemplate(table)
	for _, c := range columns {
		b.columns = append(b.columns, NewTemplate(c))
	}
	return b
}

func (b *Builder) SelectModel(model interface{}, ignoreColumns ...string) *Builder {
	if b.err != nil {
		return b
	}
	t, err := ParseTable(b.dialect, model)
	if err != nil {
		b.err = err
		return b
	}
	names, err := t.ColumnNames()
	if err != nil {
		b.err = err
		return b
	}
	b.action = actionSelect
	b.table = NewTemplate(ParseTableName(model))
	for _, name := range names {
		if containStrings(ignoreColumns, name) {
			continue
		}
		b.columns = append(b.columns, NewTemplate(name))
	}
	return b
}

func (b *Builder) Insert(table string, columns map[string]interface{}) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionInsert
	b.table = NewTemplate(table)
	for name, value := range columns {
		b.columns = append(b.columns, NewTemplate(name, value))
	}
	return b
}

func (b *Builder) InsertModel(model interface{}, ignoreZeroValue bool, ignoreColumns ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionInsert
	b.table = NewTemplate(ParseTableName(model))
	pairs, err := ParseColumnValuePairs(b.dialect, model)
	if err != nil {
		b.err = err
		return b
	}
	for _, pair := range pairs {
		if ignoreZeroValue && reflect.ValueOf(pair.Value).IsZero() {
			continue
		}
		if containStrings(ignoreColumns, pair.Column) {
			continue
		}
		b.columns = append(b.columns, NewTemplate(pair.Column, pair.Value))
	}
	return b
}

func (b *Builder) Update(table string, columns map[string]interface{}) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionUpdate
	b.table = NewTemplate(table)
	for name, value := range columns {
		b.columns = append(b.columns, NewTemplate(name, value))
	}
	return b
}

func (b *Builder) UpdateModel(model interface{}, ignoreZeroValue bool, ignoreFields ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionUpdate
	b.table = NewTemplate(ParseTableName(model))
	m, err := mapStructToMap(model)
	if err != nil {
		b.err = err
		return b
	}
	for name, value := range m {
		if ignoreZeroValue && reflect.ValueOf(value).IsZero() {
			continue
		}
		if containStrings(ignoreFields, name) {
			continue
		}
		b.columns = append(b.columns, NewTemplate(name, value))
	}
	return b
}

func (b *Builder) Delete(table string) *Builder {
	if b.err != nil {
		return b
	}
	b.action = actionDelete
	b.table = NewTemplate(table)
	return b
}

func (b *Builder) Where(format string, values ...interface{}) *Builder {
	if b.err != nil {
		return b
	}
	b.where.Appendf(format, values...)
	return b
}

func (b *Builder) WhereIn(column string, values ...interface{}) *Builder {
	if b.err != nil {
		return b
	}
	if len(values) == 0 {
		values = append(values, "null")
	}
	holders := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		holders = append(holders, "?")
	}
	format := fmt.Sprintf("%s in (%s)", b.dialect.Quote(column), strings.Join(holders, ", "))
	b.where.Appendf(format, values...)
	return b
}

func (b *Builder) WhereTemplate(tt ...Template) *Builder {
	if b.err != nil {
		return b
	}
	b.where.Append(tt...)
	return b
}

func (b *Builder) OrderBy(fields ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.orderBy = append(b.orderBy, fields...)
	return b
}

func (b *Builder) Paging(page, size int) *Builder {
	if b.err != nil {
		return b
	}
	b.paging = NewTemplate("limit ?,?", (page-1)*size, size)
	return b
}

func (b *Builder) GroupBy(fields ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.groupBy = append(b.groupBy, fields...)
	return b
}

func (b *Builder) Having(format string, values ...interface{}) *Builder {
	if b.err != nil {
		return b
	}
	b.having.Appendf(format, values...)
	return b
}
