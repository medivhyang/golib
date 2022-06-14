package errors

import (
	"errors"
	"testing"
)

func TestStackError(t *testing.T) {
	var stackError ErrorStack
	t.Log(stackError.Pop())
	stackError.Push(errors.New("err1"))
	stackError.Push(errors.New("err2"))
	t.Log(stackError.Error())
	t.Log("---")
	t.Log(stackError.Top())
	t.Log(stackError.Pop())
	t.Log(stackError.Top())
	t.Log(stackError.IsEmpty())
	t.Log(stackError.Pop())
	t.Log(stackError.IsEmpty())
	t.Log(stackError)
}
