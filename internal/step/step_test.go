package step

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipStep(t *testing.T) {
	reason := "this is a test"
	err := Skip(reason)
	require.Error(t, err)
	require.Equal(t, reason, err.Error())
}

func TestSkipf(t *testing.T) {
	require.True(t, IsSkip(Skipf("whatever %s", "blah")))
}

func TestIsSkip(t *testing.T) {
	require.True(t, IsSkip(Skip("whatever")))
	require.False(t, IsSkip(errors.New("nope")))
}
