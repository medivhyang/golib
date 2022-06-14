package log

import (
	"os"
)

type ConsoleAppender struct {
	Formatter Formatter
	Filters   []Filter
}

func NewConsoleAppender(formatter Formatter, filters ...Filter) *ConsoleAppender {
	return &ConsoleAppender{Formatter: formatter, Filters: filters}
}

func (a *ConsoleAppender) Append(t Template) error {
	for _, filter := range a.Filters {
		if filter == nil {
			continue
		}
		if ok := filter(t); !ok {
			return nil
		}
	}
	if a.Formatter == nil {
		return nil
	}
	bs, err := a.Formatter.Format(t)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(bs); err != nil {
		return err
	}
	return nil
}
