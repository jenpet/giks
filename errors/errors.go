package errors

import (
	"errors"
	"fmt"
	"reflect"
)

type WarningError struct {
	err error
}

func (ee WarningError) Error() string {
	return ee.err.Error()
}

func NewWarningError(msg string) error {
	return WarningError{
		errors.New(msg),
	}
}

func NewWarningErrorf(format string, args ...interface{}) error {
	return NewWarningError(fmt.Sprintf(format, args...))
}

func New(msg string) error {
	return errors.New(msg)
}

func IsWarningError(err error) bool {
	if t := reflect.TypeOf(err); t != nil {
		return t.Name() == "WarningError"
	}
	return false
}
