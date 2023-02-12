package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoalesce(t *testing.T) {
	runCoalesceTest(t, Ptr(1), nil, nil, Ptr(1))
	runCoalesceTest[int](t, nil, nil, nil, nil)
	runCoalesceTest(t, Ptr(1), nil, nil, Ptr(1), nil)
	runCoalesceTest(t, Ptr(1), Ptr(1))
	runCoalesceTest(t, Ptr(1), Ptr(1), nil, nil)
	runCoalesceTest(t, Ptr(false), nil, nil, Ptr(false))
	runCoalesceTest(t, Ptr(false), Ptr(false))
}

func runCoalesceTest[T comparable](t *testing.T, expected *T, values ...*T) {
	result := Coalesce(values...)

	if expected == nil {
		require.Nil(t, result)
		return
	}

	require.Equal(t, *expected, *result)
}

func TestCoalesceWithFunc(t *testing.T) {
	tables := []struct {
		matchExpected bool
		expectedVal   int
		values        []*int
	}{
		{true, 6, []*int{nil, nil, Ptr(1), nil, Ptr(4), Ptr(6), nil, Ptr(8)}},
		{false, 0, []*int{nil, nil, Ptr(1), nil, Ptr(4), Ptr(3), nil, Ptr(5)}},
		{false, 0, []*int{nil, nil, nil}},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
			result, ok := CoalesceWithFunc(func(i *int) bool {
				return *i > 5
			}, table.values...)

			require.Equal(t, table.matchExpected, ok)

			if !table.matchExpected {
				require.Nil(t, result)
				return
			}

			require.NotNil(t, result, "Result should not have been nil but was")
			require.Equal(t, table.expectedVal, *result)
		})
	}
}

func TestMustCoalesce(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		expected := Ptr(1)

		require.Equal(t, expected, MustCoalesce(nil, nil, nil, expected))
	})

	t.Run("all nils", func(t *testing.T) {
		require.Panics(t, func() {
			_ = MustCoalesce[int](nil, nil, nil, nil)
		})
	})

}
