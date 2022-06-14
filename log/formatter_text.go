package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type TextFormatter struct {
	TimeLayout string
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(t Template) ([]byte, error) {
	buf := bytes.Buffer{}

	finalTimeLayout := f.TimeLayout
	if finalTimeLayout == "" {
		finalTimeLayout = time.RFC3339
	}
	buf.WriteString(t.Time.Format(finalTimeLayout))
	if LevelText(t.Level) != "" {
		buf.WriteString(fmt.Sprintf(": [%s]", LevelText(t.Level)))
	}
	if t.Prefix != "" {
		buf.WriteString(fmt.Sprintf(": [%s]", t.Prefix))
	}
	if len(t.Modules) > 0 {
		for _, item := range t.Modules {
			buf.WriteString(fmt.Sprintf(": [%s]", item))
		}
	}
	buf.WriteString(fmt.Sprintf(": %s", t.Message))
	if t.Data != nil {
		bs, err := json.Marshal(t.Data)
		if err != nil {
			buf.WriteString(fmt.Sprintf(": <json marshal fail: %s>", err.Error()))
		} else {
			buf.WriteString(fmt.Sprintf(": %s", string(bs)))
		}
	}
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
