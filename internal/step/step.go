package step

import "errors"

func IsSkip(err error) bool {
	return errors.As(err, &ErrSkip{})
}

type ErrSkip struct {
	reason string
}

// Error implements the error interface. returns the reason the pipe was skipped.
func (e ErrSkip) Error() string {
	return e.reason
}

// Skip skips this pipe with the given reason.
func Skip(reason string) ErrSkip {
	return ErrSkip{reason: reason}
}
