package orm

import (
	"fmt"
	"strings"
)

func CreateTables(dialect Dialect, models ...interface{}) TemplateWithError {
	return createTables(dialect, false, models...)
}

func CreateTablesIfNotExists(dialect Dialect, models ...interface{}) TemplateWithError {
	return createTables(dialect, true, models...)
}

func createTables(dialect Dialect, checkExists bool, models ...interface{}) TemplateWithError {
	if dialect == nil {
		dialect = GetDefaultDialect()
	}

	tt, err := ParseTables(dialect, models...)
	if err != nil {
		return TemplateWithError{Err: err}
	}

	b := strings.Builder{}
	for i, t := range tt {
		if len(t.Columns) == 0 {
			return TemplateWithError{}
		}
		if checkExists {
			b.WriteString(fmt.Sprintf("create table if not exists %s (", t.Name))
		} else {
			b.WriteString(fmt.Sprintf("create table %s (", t.Name))
		}
		for j, c := range t.Columns {
			b.WriteString(fmt.Sprintf("%s %s", dialect.Quote(c.Name), c.Type))
			suffix := strings.TrimSpace(c.Suffix)
			if suffix != "" {
				b.WriteString(" ")
				b.WriteString(suffix)
			}
			if j < len(t.Columns)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(");")
		if i < len(tt)-1 {
			b.WriteString(" ")
		}
	}

	return TemplateWithError{Template: NewTemplate(b.String())}
}

func DropTables(dialect Dialect, models ...interface{}) TemplateWithError {
	return dropTables(dialect, false, models...)
}

func DropTablesIfExists(dialect Dialect, models ...interface{}) TemplateWithError {
	return dropTables(dialect, true, models...)
}

func dropTables(dialect Dialect, checkExists bool, models ...interface{}) TemplateWithError {
	if dialect == nil {
		dialect = GetDefaultDialect()
	}

	names := make([]string, 0, len(models))
	for _, m := range models {
		names = append(names, ParseTableName(m))
	}

	b := strings.Builder{}
	for i, name := range names {
		if checkExists {
			b.WriteString(fmt.Sprintf("drop table if exists %s;", name))
		} else {
			b.WriteString(fmt.Sprintf("drop table %s;", name))
		}
		if i < len(names)-1 {
			b.WriteString(" ")
		}
	}

	return TemplateWithError{Template: NewTemplate(b.String())}
}
