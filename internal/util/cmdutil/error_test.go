package cmdutil_test

import (
	"errors"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/stretchr/testify/require"
)

func TestExitError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		const errMsg = "some error exit thinger"
		err := errors.New(errMsg)
		exitErr := cmdutil.WrapError(err, "oh no!")

		require.IsType(t, &cmdutil.ExitError{}, exitErr)
		require.ErrorContains(t, exitErr, errMsg)

	})
}
