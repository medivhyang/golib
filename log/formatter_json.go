package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type JSONFormatter struct {
	Pretty     bool
	Prefix     string
	Indent     string
	TimeLayout string
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func NewJSONFormatterPretty(prefix string, indent string) *JSONFormatter {
	return &JSONFormatter{
		Pretty: true,
		Prefix: prefix,
		Indent: indent,
	}
}

func (f *JSONFormatter) Format(t Template) ([]byte, error) {
	finalTimeLayout := f.TimeLayout
	if finalTimeLayout == "" {
		finalTimeLayout = time.RFC3339
	}

	view := struct {
		Time    string                 `json:"time,omitempty"`
		Level   string                 `json:"level,omitempty"`
		Prefix  string                 `json:"prefix,omitempty"`
		Module  []string               `json:"module,omitempty"`
		Message string                 `json:"message,omitempty"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}{
		Prefix:  t.Prefix,
		Module:  t.Modules,
		Level:   LevelText(t.Level),
		Message: t.Message,
		Time:    t.Time.Format(finalTimeLayout),
		Data:    t.Data,
	}

	buf := bytes.Buffer{}
	if f.Pretty {
		bs, err := json.MarshalIndent(view, f.Prefix, f.Indent)
		if err != nil {
			buf.WriteString(fmt.Sprintf("{\"error\":\"<json marshal indent fail: %s>\"}", err.Error()))
		} else {
			buf.Write(bs)
		}
	} else {
		bs, err := json.Marshal(view)
		if err != nil {
			buf.WriteString(fmt.Sprintf("{\"error\":\"<json marshal fail: %s>\"}", err.Error()))
		} else {
			buf.Write(bs)
		}
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}
