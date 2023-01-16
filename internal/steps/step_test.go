package steps

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestStep_FindAllChildren(t *testing.T) {
	project, _ := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
	step := DeploymentStep(project)

	t.Run("list all children", func(t *testing.T) {
		children := step.FindAllChildren(nil)
		for _, child := range children {
			require.NotNil(t, child)
			require.IsType(t, &Step{}, child)
		}
	})

	t.Run("list specific type", func(t *testing.T) {
		children := step.FindAllChildren(aws.String("TaskDefinition"))

		for _, child := range children {
			require.NotNil(t, child)
			require.IsType(t, &Step{}, child)
			require.Equal(t, "TaskDefinition", child.Label)
		}

	})
}
