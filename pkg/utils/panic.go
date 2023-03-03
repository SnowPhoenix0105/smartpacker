package utils

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type PanicError struct {
	Msg   any
	Stack []byte
}

func NewPanicError(msg any) error {
	return &PanicError{Msg: msg, Stack: debug.Stack()}
}

func (pe *PanicError) Error() string {
	return fmt.Sprintf("panic [recovered]: %v\n%s", pe.Msg, string(pe.Stack))
}

func (pe *PanicError) Unwrap() error {
	if err, ok := pe.Msg.(error); ok {
		return err
	}
	return nil
}

func PanicRecoverWithFinalErrorPtr(finalErrPtr *error) {
	msg := recover()
	if msg == nil {
		return
	}

	panicErr := NewPanicError(msg)
	if finalErrPtr == nil {
		println(panicErr.Error())
	}

	if oldErr := *finalErrPtr; oldErr != nil {
		*finalErrPtr = errors.Join(panicErr, oldErr)
		return
	}

	*finalErrPtr = panicErr
}
