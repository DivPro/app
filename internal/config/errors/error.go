package errors

import (
	"fmt"
	"reflect"
)

type InvalidError struct {
	Err   error
	Cause string
}

func (err InvalidError) Error() string {
	return fmt.Sprintf("invalid [%s]: %s", err.Cause, err.Err.Error())
}

type ValidationError struct {
	Subject any
	Errors  []InvalidError
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("validation of %T failed: %v",
		reflect.Indirect(reflect.ValueOf(err.Subject)).Interface(),
		err.Errors,
	)
}
