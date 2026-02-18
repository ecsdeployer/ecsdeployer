package cmdutil

type ExitError struct {
	Err     error
	Code    int
	Details string
}

func WrapErrorWithCode(err error, code int, details string) *ExitError {
	return &ExitError{
		Err:     err,
		Code:    code,
		Details: details,
	}
}

func WrapError(err error, log string) *ExitError {
	return WrapErrorWithCode(err, 1, log)
}

func (e *ExitError) Error() string {
	return e.Err.Error()
}
