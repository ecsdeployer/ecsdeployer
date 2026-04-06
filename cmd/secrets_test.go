package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/usererr"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestSecretsDelete(t *testing.T) {
	helpers.IsTestingMode = true
	t.Run("ok", func(t *testing.T) {

		testutil.StartMocker(t, awsmocker.WithMocks(
			testutil.Mock_SSM_DeleteParameter("/ecsdeployer/dummy/DUMMY_VAR_TO_REMOVE"),
		))

		result := runCommand(t, nil, "secrets", "delete", "-c", "../internal/builders/testdata/everything.yml", "DUMMY_VAR_TO_REMOVE")
		require.NoError(t, result.err)
	})

	t.Run("notfound", func(t *testing.T) {

		testutil.StartMocker(t, awsmocker.WithMocks(
			testutil.Mock_SSM_DeleteParameter_NotFound("/ecsdeployer/dummy/DUMMY_VAR_TO_REMOVE"),
		))

		result := runCommand(t, nil, "secrets", "delete", "-c", "../internal/builders/testdata/everything.yml", "DUMMY_VAR_TO_REMOVE")
		require.Error(t, result.err)
		require.ErrorAs(t, result.err, new(*usererr.UserError))
		require.Contains(t, result.stderr, "ParameterNotFound")
	})

	t.Run("forbidden", func(t *testing.T) {

		testutil.StartMocker(t, awsmocker.WithMocks(
			testutil.Mock_SSM_DeleteParameter_Forbidden("/ecsdeployer/dummy/DUMMY_VAR_TO_REMOVE"),
		))

		result := runCommand(t, nil, "secrets", "delete", "-c", "../internal/builders/testdata/everything.yml", "DUMMY_VAR_TO_REMOVE")
		require.Error(t, result.err)
		require.ErrorAs(t, result.err, new(*usererr.UserError))
		require.Contains(t, result.stderr, "AccessDenied")
	})
}
