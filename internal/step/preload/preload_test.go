package preload

import (
	"errors"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type fakePreloader struct {
	shouldSkip bool
	shouldFail bool
	failErr    error
	runCount   int
}

func (fakePreloader) String() string { return "fake" }
func (f *fakePreloader) Preload(ctx *config.Context) error {
	f.runCount++
	switch {
	case f.shouldSkip:
		return step.Skip("step skipped")
	case f.shouldFail:
		if f.failErr != nil {
			return f.failErr
		}
		return errors.New("some failure")
	default:
		return nil
	}
}

func TestPreloadStep(t *testing.T) {
	testutil.DisableLoggingForTest(t)
	fakeCtx := &config.Context{}

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "preloading resources", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		require.False(t, Step{}.Skip(fakeCtx))
	})

	t.Run("Run", func(t *testing.T) {

		origSteps := slices.Clone(preloaders)

		t.Cleanup(func() {
			preloaders = make([]preloader, len(origSteps))
			copy(preloaders, origSteps)
		})

		t.Run("count", func(t *testing.T) {
			require.Len(t, preloaders, 2)
		})

		t.Run("abnormal", func(t *testing.T) {

			t.Run("two skips", func(t *testing.T) {
				preloaders = []preloader{
					&fakePreloader{shouldSkip: true},
					&fakePreloader{shouldSkip: true},
					&fakePreloader{shouldSkip: true},
				}
				require.NoError(t, Step{}.Run(fakeCtx))
			})

			t.Run("error", func(t *testing.T) {
				preloaders = []preloader{
					&fakePreloader{shouldSkip: true},
					&fakePreloader{shouldFail: true},
				}
				require.Error(t, Step{}.Run(fakeCtx))
			})
		})

		t.Run("success", func(t *testing.T) {
			fakePl := &fakePreloader{}
			preloaders = []preloader{fakePl}
			require.NoError(t, Step{}.Run(fakeCtx))
			require.Equal(t, 1, fakePl.runCount)
		})

	})
}
