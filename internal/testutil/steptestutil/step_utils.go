package steptestutil

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func StepTestAwsMocker(t *testing.T, projectFilePath string, mocks []*awsmocker.MockedEndpoint) (*config.Project, *config.Context) {
	helpers.IsTestingMode = true

	testutil.DisableLoggingForTest(t)

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: append([]*awsmocker.MockedEndpoint{}, mocks...),
	})

	project, err := yaml.ParseYAMLFile[config.Project](projectFilePath)
	require.NoError(t, err)

	ctx := config.New(project)

	return project, ctx
}
