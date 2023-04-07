package cmd

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	dsmock "ecsdeployer.com/ecsdeployer/internal/testutil/mocks/ecs/describeservicemock"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanCmd(t *testing.T) {
	silenceLogging(t)
	helpers.IsTestingMode = true

	t.Run("calls correct function", func(t *testing.T) {

		testutil.StartMocker(t, &awsmocker.MockerOptions{
			Mocks: []*awsmocker.MockedEndpoint{
				testutil.Mock_Tagging_GetResources("ecs:task-definition", map[string]string{"ecsdeployer/project": "dummy/fancy"}, []string{
					fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-fancy-svc1:555", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
					fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-fancy-old-service:122", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
					fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-fancy-old-service:123", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
					fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-fancy-other-service:123", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
				}),

				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-fancy-old-service", 122),
				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-fancy-old-service", 123),
				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-fancy-other-service", 123),

				testutil.Mock_Tagging_GetResources("ecs:service", map[string]string{"ecsdeployer/project": "dummy/fancy"}, []string{
					fmt.Sprintf("arn:aws:ecs:%s:%s:service/dummy/dummy-fancy-svc1", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
					fmt.Sprintf("arn:aws:ecs:%s:%s:service/dummy/dummy-fancy-old-service", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
					fmt.Sprintf("arn:aws:ecs:%s:%s:service/dummy/dummy-fancy-other-service", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
				}),
				dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-fancy-old-service")),
				dsmock.Mock(dsmock.WithStable(), dsmock.WithName("dummy-fancy-other-service")),
				testutil.Mock_ECS_UpdateService_Generic(),
				testutil.Mock_ECS_DeleteService_jmespath(map[string]any{"service": "dummy-fancy-old-service"}, ecsTypes.Service{}),
				testutil.Mock_ECS_DeleteService_jmespath(map[string]any{"service": "dummy-fancy-other-service"}, ecsTypes.Service{}),

				testutil.Mock_Scheduler_GetScheduleGroup("dummy-fancy"),
				testutil.Mock_Scheduler_ListSchedules("dummy-fancy", []testutil.MockListScheduleEntry{
					{Name: "ecsd-cron-dummy-fancy-cron1"},
					{Name: "ecsd-cron-dummy-fancy-cron-old"},
				}),
				testutil.Mock_Scheduler_DeleteSchedule("dummy-fancy", "ecsd-cron-dummy-fancy-cron-old"),
			},
		})

		result := runCommand(t, "clean", "-q", "-c", "testdata/info_simple.yml")

		require.NoError(t, result.err)
		require.Equal(t, 0, result.exitCode)

		// var checkErr *exitError
		// require.ErrorAs(t, result.err, &checkErr)
	})
}
