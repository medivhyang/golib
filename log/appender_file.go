package log

import (
	"os"
)

type FileAppender struct {
	File      *os.File
	Formatter Formatter
	Filters   []Filter
}

func NewFileAppender(filePath string, formatter Formatter, filters ...Filter) (*FileAppender, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0744)
	if err != nil {
		return nil, err
	}
	return &FileAppender{
		Formatter: formatter,
		File:      file,
		Filters:   filters,
	}, nil
}

func (a *FileAppender) Append(t Template) error {
	for _, filter := range a.Filters {
		if filter == nil {
			continue
		}
		if ok := filter(t); !ok {
			return nil
		}
	}
	if a.File == nil {
		return nil
	}
	if a.Formatter == nil {
		return nil
	}
	bs, err := a.Formatter.Format(t)
	if err != nil {
		return err
	}
	if _, err := a.File.Write(bs); err != nil {
		return err
	}
	return nil
}
