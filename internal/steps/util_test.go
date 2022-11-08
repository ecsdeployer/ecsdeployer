package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func stepTestAwsMocker(t *testing.T, projectFilePath string, mocks []*awsmocker.MockedEndpoint) (*config.Project, *config.Context) {
	awsmocker.Start(t, &awsmocker.MockerOptions{
		Mocks: append([]*awsmocker.MockedEndpoint{}, mocks...),
	})

	project, err := yaml.ParseYAMLFile[config.Project](projectFilePath)
	require.NoError(t, err)

	ctx := config.New(project)

	return project, ctx
}
