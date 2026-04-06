package cmdutil

import (
	"context"
	"errors"
	"fmt"
)

// some of this is copied from GithubCLI

// Error unrelated to usage
// FlagErrorf returns a new FlagError that wraps an error produced by
// fmt.Errorf(format, args...).
func FlagErrorf(format string, args ...any) error {
	return FlagErrorWrap(fmt.Errorf(format, args...))
}

// FlagErrorWrap returns a new FlagError that wraps the specified error.
func FlagErrorWrap(err error) error { return &FlagError{err} }

// A *FlagError indicates an error processing command-line flags or other arguments.
// Such errors cause the application to display the usage message.
type FlagError struct {
	// Note: not struct{error}: only *FlagError should satisfy error.
	err error
}

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

func (fe *FlagError) Unwrap() error {
	return fe.err
}

// ErrSilentError is an error that triggers exit code 1 without any error messaging
var ErrSilentError = errors.New("SilentError")

// ErrCancelError signals user-initiated cancellation
var ErrCancelError = errors.New("CancelError")

// ErrPendingError signals nothing failed but something is pending
var ErrPendingError = errors.New("PendingError")

func IsUserCancellation(err error) bool {
	return errors.Is(err, ErrCancelError) || errors.Is(err, context.Canceled)
}
