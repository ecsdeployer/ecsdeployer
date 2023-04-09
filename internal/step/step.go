package step

import (
	"errors"
	"fmt"
)

func IsSkip(err error) bool {
	return errors.As(err, &ErrSkip{})
}

type ErrSkip struct {
	reason string
}

// Error implements the error interface. returns the reason the step was skipped.
func (e ErrSkip) Error() string {
	return e.reason
}

// Skip skips this step with the given reason.
func Skip(reason string) ErrSkip {
	return ErrSkip{reason: reason}
}

func Skipf(msg string, fields ...any) ErrSkip {
	return Skip(fmt.Sprintf(msg, fields...))
}
