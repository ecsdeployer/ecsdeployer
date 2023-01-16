package testutil

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/stretchr/testify/require"
)

// Does one object serialize the same as another object? (generate the same JSON)
// func ObjectSerializationMatches(t *testing.T, expected any, actual any) bool {
// 	t.Helper()
// 	actualJson, err := util.Jsonify(actual)
// 	require.NoError(t, err, "unable to serialize <actual>")

// 	expJson, err := util.Jsonify(expected)
// 	require.NoError(t, err, "Unable to serialize <expected>")

// 	return actualJson == expJson
// }

func RequireObjectSerializationEqual(t *testing.T, expected any, actual any) {
	t.Helper()
	actualJson, err := util.Jsonify(actual)
	require.NoError(t, err, "unable to serialize <actual>")

	expJson, err := util.Jsonify(expected)
	require.NoError(t, err, "Unable to serialize <expected>")

	require.Equal(t, expJson, actualJson)
}
