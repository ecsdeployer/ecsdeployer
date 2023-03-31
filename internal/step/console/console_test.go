package console

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestConsoleStep(t *testing.T) {

	testutil.DisableLoggingForTest(t)

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "registering console task", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		t.Run("console disabled", func(t *testing.T) {
			project := &config.Project{
				ConsoleTask: &config.ConsoleTask{
					Enabled: aws.Bool(false),
				},
			}
			project.ApplyDefaults()
			project.ConsoleTask.ApplyDefaults()
			ctx := config.New(project)
			require.True(t, Step{}.Skip(ctx))
		})

		t.Run("not disabled", func(t *testing.T) {
			project := &config.Project{
				ConsoleTask: &config.ConsoleTask{
					Enabled: aws.Bool(true),
				},
			}
			project.ApplyDefaults()
			project.ConsoleTask.ApplyDefaults()
			ctx := config.New(project)
			require.False(t, Step{}.Skip(ctx))
		})
	})

	t.Run("Run", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Logs_CreateLogGroup_AllowAny(),
				testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
				testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
			})
			err := Step{}.Run(ctx)
			require.NoError(t, err)
		})
	})
}
