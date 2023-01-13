package config

import (
	"errors"
	"fmt"
)

var ErrValidation = errors.New("config validation error")

type ValidationError struct {
	Reason string
	Err    error
}

// ensure adherence to interface
var _ error = (*ValidationError)(nil)

func (ve *ValidationError) Error() string {
	if ve.Reason != "" {
		return ve.Reason
	}

	if ve.Err != nil {
		return ve.Err.Error()
	}

	return ErrValidation.Error()
}

func (ve *ValidationError) Unwrap() error {
	return ve.Err
}

func (ve *ValidationError) Is(err error) bool {
	if err == nil {
		return false
	}

	//nolint:errorlint
	if err == ErrValidation {
		return true
	}

	if ve.Err != nil {
		return errors.Is(ve.Err, err)
	}

	return false
}

// Builds a new validation error
// If first parameter is an 'error', then validation error will wrap that
// First parameter is a string with optional format, followed by values to be formatted
// the last value may be an error object. if present, that will be wrapped.
// do not include it in your format string
func NewValidationError(value any, values ...any) *ValidationError {

	switch firstVal := value.(type) {
	case string:
		return &ValidationError{
			Reason: fmt.Sprintf(firstVal, values...),
		}

	case error:
		if len(values) > 0 {
			fmtString, rest := values[0], values[1:]

			return &ValidationError{
				Reason: fmt.Sprintf(fmtString.(string), rest...),
				Err:    firstVal,
			}
		} else {
			return &ValidationError{
				Reason: firstVal.Error(),
				Err:    firstVal,
			}
		}

	default:
		return &ValidationError{Reason: "Validation Error"}
	}
}
