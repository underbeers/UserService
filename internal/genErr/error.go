package genErr

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Error struct {
	prevError error
	message   string
	data      map[string]any
}

func (se *Error) Error() string {
	return fmt.Sprintf("%s", se.message) //nolint: gosimple
}

func (se *Error) ErrorFull() string {
	if len(se.data) == 0 {
		return fmt.Sprintf("%s; ", se.message)
	}

	return fmt.Sprintf("%s: %v; ", se.message, se.data)
}

func (se *Error) ErrorEx() string {
	var str string
	var err error = se

	for {
		e, ok := err.(interface { //nolint: errorlint
			ErrorFull() string
		})
		if !ok {
			str += err.Error() + "; "
		} else {
			str += e.ErrorFull()
		}

		u, ok := err.(interface { //nolint: errorlint
			Unwrap() error
		})
		if !ok {
			return str
		}
		err = u.Unwrap()
		if err == nil {
			return str
		}
	}
}

func NewError(prevError error, message error, data ...any) *Error {
	d := make(map[string]any)
	var k string

	if len(data)%2 != 0 {
		data = data[:len(data)-1] //nolint: staticcheck
	} else {
		for i, v := range data {
			if i%2 == 0 {
				k = fmt.Sprintf("%v", v)
			} else {
				d[k] = v
			}
		}
	}

	se := &Error{
		prevError: prevError,
		message:   message.Error(),
		data:      d,
	}

	if prevError != nil && message != nil {
		pe, pOk := prevError.(*Error) //nolint: errorlint

		if pOk {
			if reflect.DeepEqual(pe.data, se.data) && pe.message == se.message {
				return pe
			}
		}
	}

	return se
}

func (se *Error) As(target any) bool {
	var msg string
	switch v := target.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		log.Fatal("genError target any param unknown type")
	}

	return se.message == msg
}

func (se *Error) Is(target error) bool {
	t := target.Error()

	return se.message == t
}

func (se *Error) Unwrap() error {
	return se.prevError
}

func New(text string) *Error {
	return NewError(nil, errors.New(text)) //nolint: goerr113
}
