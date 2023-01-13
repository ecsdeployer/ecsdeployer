package config_test

import (
	"fmt"
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {

	t.Run("Is_Unraps", func(t *testing.T) {

		tables := []struct {
			err     *config.ValidationError
			is      []error
			isnot   []error
			unwraps error
		}{
			{config.NewValidationError(os.ErrClosed), []error{os.ErrClosed}, nil, os.ErrClosed},
			{config.NewValidationError("something"), nil, nil, nil},
		}

		for tNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", tNum), func(t *testing.T) {
				require.False(t, table.err.Is(nil))
				require.True(t, table.err.Is(config.ErrValidation))

				require.NotErrorIs(t, table.err, nil)
				require.ErrorIs(t, table.err, config.ErrValidation)

				if table.unwraps != nil {
					require.ErrorIs(t, table.err.Unwrap(), table.unwraps)
				} else {
					require.Nil(t, table.err.Unwrap())
				}

				if len(table.is) > 0 {
					for _, serr := range table.is {
						require.ErrorIs(t, table.err, serr)
					}
				}

				if len(table.isnot) > 0 {
					for _, serr := range table.isnot {
						require.NotErrorIs(t, table.err, serr)
					}
				}
			})
		}
	})

	t.Run("NewValidationError", func(t *testing.T) {
		tables := []struct {
			err     *config.ValidationError
			msg     string
			wrapped error
		}{
			{config.NewValidationError("some error"), "some error", nil},
			{config.NewValidationError(os.ErrClosed), "file already closed", os.ErrClosed},
			{config.NewValidationError(os.ErrClosed, "some error"), "some error", os.ErrClosed},
			{config.NewValidationError(os.ErrClosed, "test %d %s %d", 1, "xx", 5), "test 1 xx 5", os.ErrClosed},
			{config.NewValidationError("test %d %s %d", 1, "xx", 5), "test 1 xx 5", nil},

			// edge
			{&config.ValidationError{Err: os.ErrClosed}, "file already closed", os.ErrClosed},
			{&config.ValidationError{}, config.ErrValidation.Error(), nil},

			// Lazy dev cases
			{config.NewValidationError(nil), "Validation Error", nil},
			{config.NewValidationError(1), "Validation Error", nil},
			{config.NewValidationError(1.0), "Validation Error", nil},
			{config.NewValidationError(false), "Validation Error", nil},
			{config.NewValidationError(true), "Validation Error", nil},
		}

		for tNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", tNum), func(t *testing.T) {
				require.ErrorIs(t, table.err, config.ErrValidation)

				if table.msg != "" {
					require.ErrorContains(t, table.err, table.msg)
				}

				if table.wrapped != nil {
					require.ErrorIs(t, table.err, table.wrapped)
					require.Equal(t, table.wrapped, table.err.Unwrap())
				} else {
					require.Nil(t, table.err.Unwrap())
				}
			})
		}
	})
}
