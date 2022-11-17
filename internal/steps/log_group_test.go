package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestLogGroupStep(t *testing.T) {
	t.Run("when logging disabled globally", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
		ctx.Project.Logging.Disabled = true

		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when logging disabled for specific task", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})

		project.Services[0].LoggingConfig = &config.TaskLoggingConfig{
			Driver: aws.String(config.LoggingDisableFlag),
		}

		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group not created, unlimited retention", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{}),
			testutil.Mock_Logs_CreateLogGroup("/ecsdeployer/app/dummy/web"),
			// testutil.Mock_Logs_PutRetentionPolicy("/ecsdeployer/app/dummy/web", 30),
		})
		retObj, _ := config.ParseLogRetention("forever")
		project.Logging.AwsLogConfig.Retention = &retObj

		// fake the preload step as it will always run before
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))

		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group not created, limit retention", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{}),
			testutil.Mock_Logs_CreateLogGroup("/ecsdeployer/app/dummy/web"),
			testutil.Mock_Logs_PutRetentionPolicy("/ecsdeployer/app/dummy/web", 30),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Logging.AwsLogConfig.Retention = &retObj

		// fake the preload step as it will always run before
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))

		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group exists, has correct retention", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{
				"/ecsdeployer/app/dummy/web": 30,
			}),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Logging.AwsLogConfig.Retention = &retObj

		// fake the preload step as it will always run before
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))

		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group exists, ret to unlimited", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{
				"/ecsdeployer/app/dummy/web": 30,
			}),
			testutil.Mock_Logs_DeleteRetentionPolicy("/ecsdeployer/app/dummy/web"),
		})
		retObj, _ := config.ParseLogRetention("forever")
		project.Logging.AwsLogConfig.Retention = &retObj
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))
		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group exists, ret to diff ret", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{
				"/ecsdeployer/app/dummy/web": 60,
			}),
			testutil.Mock_Logs_PutRetentionPolicy("/ecsdeployer/app/dummy/web", 30),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Logging.AwsLogConfig.Retention = &retObj
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))
		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group exists, unlimit to ret", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{
				"/ecsdeployer/app/dummy/web": 0,
			}),
			testutil.Mock_Logs_PutRetentionPolicy("/ecsdeployer/app/dummy/web", 30),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Logging.AwsLogConfig.Retention = &retObj
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))
		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group exists, ret mismatch, sync disabled", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{
				"/ecsdeployer/app/dummy/web": 0,
			}),
			testutil.Mock_Logs_PutRetentionPolicy("/ecsdeployer/app/dummy/web", 30),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Settings.KeepInSync.LogRetention = aws.Bool(false)
		project.Logging.AwsLogConfig.Retention = &retObj
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))
		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("log group create, ret error", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{}),
			testutil.Mock_Logs_CreateLogGroup("/ecsdeployer/app/dummy/web"),
			awsmocker.Mock_Failure("logs", "PutRetentionPolicy"),
		})
		retObj, _ := config.ParseLogRetention(30)
		project.Logging.AwsLogConfig.Retention = &retObj
		require.NoError(t, PreloadLogGroupsStep(project).Apply(ctx))
		err := LogGroupStep(project.Services[0]).Apply(ctx)
		require.Error(t, err)
	})

}
