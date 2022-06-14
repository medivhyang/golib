package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type Condition struct {
	dialect   Dialect
	templates []Template
	err       error
}

func NewCondition(dialect Dialect, templates ...Template) *Condition {
	return &Condition{dialect: dialect, templates: templates}
}

func (c *Condition) Append(others ...Template) *Condition {
	c.templates = append(c.templates, others...)
	return c
}

func (c *Condition) Appendf(format string, values ...interface{}) *Condition {
	return c.Append(NewTemplate(format, values...))
}

func (c *Condition) AppendMap(m map[string]interface{}) *Condition {
	for k, v := range m {
		c.Appendf(fmt.Sprintf("%s = ?", k), v)
	}
	return c
}

func (c *Condition) AppendStruct(i interface{}, ignoreZeroValue bool) *Condition {
	m, err := mapStructToMap(i)
	if err != nil {
		c.err = err
		return c
	}
	m2 := make(map[string]interface{}, len(m))
	for k, v := range m {
		if ignoreZeroValue && reflect.ValueOf(v).IsZero() {
			continue
		}
		m2[k] = v
	}
	return c.AppendMap(m2)
}

func (c *Condition) AppendIn(column string, values ...interface{}) *Condition {
	if len(values) == 0 {
		values = append(values, "null")
	}
	holders := repeatString("?", len(values))
	format := fmt.Sprintf("%s in (%s)", c.dialect.Quote(column), strings.Join(holders, ", "))
	return c.Appendf(format, values...)
}

func (c *Condition) Join(sep string, right, left string) Template {
	if len(c.templates) == 0 {
		return Template{}
	}
	formats := make([]string, 0, len(c.templates))
	values := make([]interface{}, 0, len(c.templates))
	for _, t := range c.templates {
		formats = append(formats, t.Format)
		values = append(values, t.Values...)
	}
	return NewTemplate(right+strings.Join(formats, sep)+left, values...)
}

func (c *Condition) And() Template {
	return c.Join(" and ", "", "")
}

func (c *Condition) Or() Template {
	return c.Join(" or ", "", "")
}

func (c *Condition) IsNotEmpty() bool {
	return len(c.templates) > 0
}

func (c *Condition) Err() error {
	return c.err
}
