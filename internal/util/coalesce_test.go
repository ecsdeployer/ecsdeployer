package util

import (
	"testing"
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
		if result != nil {
			t.Fatalf("Incorrect result for Coalesce test: expected nil but was given not nil")
		}
		return
	}

	if *expected != *result {
		t.Fatalf("Incorrect result for coalesce. expected=<%v> got=<%v / %v>", *expected, result, *result)
		return
	}
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

	for _, table := range tables {
		result, ok := CoalesceWithFunc(func(i *int) bool {
			return *i > 5
		}, table.values...)

		if ok != table.matchExpected {
			t.Fatalf("Expected to have ok=%t but got %t", table.matchExpected, ok)
		}

		if table.matchExpected {
			if result == nil {
				t.Fatalf("Result should not have been nil but was")
			}
			if table.expectedVal != *result {
				t.Fatalf("Result expected=%v got=%v", table.expectedVal, *result)
			}
		} else if result != nil {
			t.Fatalf("Result should have been nil but wasnt")
		}
	}
}

func TestMustCoalesce(t *testing.T) {

	expected := Ptr(1)

	if MustCoalesce(nil, nil, nil, expected) != expected {
		t.Fatal("MustCoalesce() is broken")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustCoalesce() should have panicked")
		}
	}()

	_ = MustCoalesce[int](nil, nil, nil, nil)
}
