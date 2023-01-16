package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonifyAndPretty(t *testing.T) {

	type jsonStruct struct {
		Thing string `json:"thing"`
	}

	tables := []struct {
		obj interface{}
		exp string
	}{
		{5, `5`},
		{"test", `"test"`},
		{true, `true`},
		{jsonStruct{Thing: "blahblahba"}, `{"thing":"blahblahba"}`},
		{&jsonStruct{Thing: "blahblahba"}, `{"thing":"blahblahba"}`},
		{[]int{1, 2, 3, 4, 5}, `[1,2,3,4,5]`},
		{nil, `null`},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			res, err := Jsonify(table.obj)
			require.NoError(t, err)
			require.JSONEq(t, table.exp, res)

			resPretty, err := JsonifyPretty(table.obj)
			require.NoError(t, err)

			require.JSONEq(t, res, resPretty)
		})
	}
}
