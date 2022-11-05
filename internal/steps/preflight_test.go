package steps

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreflightStep(t *testing.T) {
	closeFunc, project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
	defer closeFunc()

	err := PreflightStep(project).Apply(ctx)
	require.NoError(t, err)

}
