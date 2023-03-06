package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExitError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		const errMsg = "some error exit thinger"
		err := errors.New(errMsg)
		exitErr := wrapError(err, "oh no!")

		require.IsType(t, &exitError{}, exitErr)
		require.ErrorContains(t, exitErr, errMsg)

	})
}
