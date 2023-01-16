package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreloadLogGroupsStep(t *testing.T) {

	t.Run("when disabled", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_simple.yml", []*awsmocker.MockedEndpoint{})
		ctx.Project.Logging.Disabled = true

		err := PreloadLogGroupsStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when access denied", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			awsmocker.Mock_Failure("logs", "DescribeLogGroups"),
		})
		err := PreloadLogGroupsStep(project).Apply(ctx)
		require.ErrorContains(t, err, "AccessDenied")
	})

	t.Run("when no log groups", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(nil),
		})
		err := PreloadLogGroupsStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when log groups are returned", func(t *testing.T) {
		logRetentionMap := map[string]int32{
			"/ecsdeployer/app/dummy/web":       30,
			"/ecsdeployer/app/dummy/worker":    30,
			"/ecsdeployer/app/dummy/cron1":     30,
			"/ecsdeployer/app/dummy/offworker": 0,
			"/ecsdeployer/app/dummy/pd1":       30,
		}
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_DescribeLogGroups(logRetentionMap),
		})
		err := PreloadLogGroupsStep(project).Apply(ctx)
		require.NoError(t, err)

		require.Equal(t, len(logRetentionMap), len(ctx.Cache.LogGroups))

		for _, logGroup := range ctx.Cache.LogGroups {
			require.NotNil(t, logGroup)
			require.IsType(t, logTypes.LogGroup{}, logGroup)

			require.NotNil(t, logGroup.Arn)
			require.NotNil(t, logGroup.LogGroupName)

			logName := *logGroup.LogGroupName

			expectedRet := logRetentionMap[logName]

			if expectedRet > 0 {
				require.Equal(t, expectedRet, *logGroup.RetentionInDays)
			} else {
				require.Nil(t, logGroup.RetentionInDays)
			}

			logArn := *logGroup.Arn
			arn, err := arn.Parse(logArn)
			require.NoError(t, err)
			require.Equal(t, "log-group:"+logName, arn.Resource)
		}

	})
}
