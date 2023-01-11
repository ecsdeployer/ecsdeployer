package config_test

import (
	"fmt"
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {
	t.Run("Is", func(t *testing.T) {})
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
