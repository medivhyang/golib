package errors

import (
	"fmt"
	"strings"
)

type ErrorStack []error

func (s *ErrorStack) Push(err error) {
	*s = append(*s, err)
}

func (s *ErrorStack) Pop() error {
	if s.IsEmpty() {
		return nil
	}
	result := (*s)[0]
	*s = (*s)[1:]
	return result
}

func (s *ErrorStack) Top() error {
	if s.IsEmpty() {
		return nil
	}
	return (*s)[0]
}

func (s *ErrorStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *ErrorStack) Errors() []error {
	return *s
}

func (s *ErrorStack) Error() string {
	b := strings.Builder{}
	b.WriteString("stack error:\n")
	for index, item := range *s {
		b.WriteString(fmt.Sprintf("%d: ", index))
		b.WriteString(item.Error())
		if index < len(*s)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
