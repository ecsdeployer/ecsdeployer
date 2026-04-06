package usererr

import "fmt"

type UserError struct {
	error
}

var _ error = (*UserError)(nil)

func (e *UserError) Error() string {
	if e.error == nil {
		return `Unknown error occurred`
	}

	return e.error.Error()
}

func (e *UserError) Unwrap() error {
	return e.error
}

func Wrap(err error) error {
	return New(err)
}

func New(args ...any) error {

	err := &UserError{}

	if len(args) == 0 {
		return err
	}

	first := args[0]

	switch v := first.(type) {
	case string:
		err.error = fmt.Errorf(v, args[1:]...)
	case error:
		err.error = v
	default:
		err.error = fmt.Errorf(`bad dev!: %T %v`, v, args)
	}

	return err
}

func Newf(msg string, args ...any) error {
	return &UserError{
		error: fmt.Errorf(msg, args...),
	}
}
