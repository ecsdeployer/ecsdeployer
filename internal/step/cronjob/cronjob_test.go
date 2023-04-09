package cronjob

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCronjobStep(t *testing.T) {
	require.Equal(t, "cronjob", Step{}.String())
}
