package cmd

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	dsmock "ecsdeployer.com/ecsdeployer/internal/testutil/mocks/ecs/describeservicemock"
	"ecsdeployer.com/ecsdeployer/internal/testutil/mocks/ecs/taskmock"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestDeployCmd(t *testing.T) {
	helpers.IsTestingMode = true

	t.Run("failure", func(t *testing.T) {
		result := runCommand(t, nil, "deploy", "-c", "badfile.yml")
		require.Error(t, result.err)
	})

}

// This runs the whole pipeline and is mainly used for testing the output
// but it also performs a pretty thorough runthrough of the deploy pipeline
func TestDeploySmoke(t *testing.T) {

	t.Run("Trace", func(t *testing.T) {
		setupFullDeployMock(t)
		require.NoError(t, runCommand(t, &rcConf{noWrapLog: true}, "deploy", "-c", "../internal/builders/testdata/everything.yml", "--trace").err)
	})

	t.Run("Debug", func(t *testing.T) {
		setupFullDeployMock(t)
		require.NoError(t, runCommand(t, &rcConf{noWrapLog: true}, "deploy", "-c", "../internal/builders/testdata/everything.yml", "--debug").err)
	})

	t.Run("Normal", func(t *testing.T) {
		setupFullDeployMock(t)
		require.NoError(t, runCommand(t, &rcConf{noWrapLog: true}, "deploy", "-c", "../internal/builders/testdata/everything.yml").err)
	})

}

func setupFullDeployMock(t *testing.T) {
	t.Helper()
	helpers.IsTestingMode = true
	setupCmdOutput(t)

	mocks := []*awsmocker.MockedEndpoint{
		testutil.Mock_EC2_DescribeSecurityGroups_Simple(),
		testutil.Mock_EC2_DescribeSubnets_Simple(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_ELBv2_DescribeTargetGroups_Generic_Success(),
		testutil.Mock_Logs_DescribeLogGroups(map[string]int32{}),
		// testutil.Mock_Logs_CreateLogGroup("/ecsdeployer/app/dummy/console"),
		testutil.Mock_Logs_CreateLogGroup_AllowAny(),
		testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
		testutil.Mock_Scheduler_GetScheduleGroup_Missing("dummy"),
		testutil.Mock_Scheduler_CreateScheduleGroup("dummy"),
		testutil.Mock_SSM_GetParametersByPath_Advanced(func(m *testutil.Mock_ECS_GetParametersByPathOpts) {
			m.MaxCount = 2
			m.NextToken = true
			m.Path = "/ecsdeployer/dummy/"
			m.NumParams = 10
		}),
		testutil.Mock_SSM_GetParametersByPath_Advanced(func(m *testutil.Mock_ECS_GetParametersByPathOpts) {
			m.MaxCount = 1
			m.Path = "/ecsdeployer/dummy/"
			m.NumParams = 8
		}),

		testutil.Mock_Tagging_GetResources("ecs:task-definition", map[string]string{"ecsdeployer/project": "dummy"}, []string{
			fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-something-something:122", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
			fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-something-something:123", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
		}),

		testutil.Mock_Tagging_GetResources("ecs:service", map[string]string{"ecsdeployer/project": "dummy"}, []string{
			fmt.Sprintf("arn:aws:ecs:%s:%s:service/dummy/dummy-old-service", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
			fmt.Sprintf("arn:aws:ecs:%s:%s:service/dummy/dummy-other-service", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
		}),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-old-service")),
		dsmock.Mock(dsmock.WithStable(), dsmock.WithName("dummy-other-service")),
		testutil.Mock_ECS_DeleteService_jmespath(map[string]any{"service": "dummy-old-service"}, ecsTypes.Service{}),
		testutil.Mock_ECS_DeleteService_jmespath(map[string]any{"service": "dummy-other-service"}, ecsTypes.Service{}),

		testutil.Mock_ECS_DeregisterTaskDefinition("dummy-something-something", 122),
		testutil.Mock_ECS_DeregisterTaskDefinition("dummy-something-something", 123),
		dsmock.Mock(dsmock.WithMaxCount(1), dsmock.WithName("dummy-svc-sidecar-ports")),
		dsmock.Mock(dsmock.WithMaxCount(1), dsmock.WithName("dummy-svc1")),
		dsmock.Mock(dsmock.WithMaxCount(1), dsmock.WithName("dummy-svc2")),
		dsmock.Mock(dsmock.WithMaxCount(1), dsmock.WithName("dummy-svc3"), dsmock.WithMissing()),
		dsmock.Mock(dsmock.WithMaxCount(1), dsmock.WithName("dummy-svc4"), dsmock.WithMissing()),
		testutil.Mock_ECS_CreateService_Generic(),
		testutil.Mock_ECS_UpdateService_Generic(),
		dsmock.Mock(dsmock.WithMaxCount(2), dsmock.WithName("dummy-svc1"), dsmock.WithPending()),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-svc1"), dsmock.WithStable()),

		dsmock.Mock(dsmock.WithMaxCount(2), dsmock.WithName("dummy-svc2"), dsmock.WithPending()),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-svc2"), dsmock.WithStable()),

		dsmock.Mock(dsmock.WithMaxCount(2), dsmock.WithName("dummy-svc3"), dsmock.WithPending()),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-svc3"), dsmock.WithStable()),

		dsmock.Mock(dsmock.WithMaxCount(2), dsmock.WithName("dummy-svc4"), dsmock.WithPending()),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-svc4"), dsmock.WithStable()),

		dsmock.Mock(dsmock.WithMaxCount(2), dsmock.WithName("dummy-svc-sidecar-ports"), dsmock.WithPending()),
		dsmock.Mock(dsmock.WithMaxCount(0), dsmock.WithName("dummy-svc-sidecar-ports"), dsmock.WithStable()),

		testutil.Mock_Scheduler_GetSchedule_Missing("dummy", "ecsd-cron-dummy-cron1"),
		testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron2"),
		testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron-daily"),

		testutil.Mock_Scheduler_CreateSchedule("dummy", "ecsd-cron-dummy-cron1"),
		testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron2"),
		testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron-daily"),

		testutil.Mock_ECS_ListTaskDefinitions("dummy-svc-sidecar-ports", []int{997, 998, 999}),
		testutil.Mock_ECS_DeregisterTaskDefinition("dummy-svc-sidecar-ports", 997),
		testutil.Mock_ECS_DeregisterTaskDefinition("dummy-svc-sidecar-ports", 998),

		testutil.Mock_Tagging_GetResources("events:rule", map[string]string{"ecsdeployer/project": "dummy"}, []string{
			fmt.Sprintf("arn:aws:events:%s:%s:rule/dummy-rule-cron1", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
		}),

		testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{
			{Name: "ecsd-cron-dummy-cron1"},
			{Name: "ecsd-cron-dummy-cron-old"},
		}),
		testutil.Mock_Scheduler_DeleteSchedule("dummy", "ecsd-cron-dummy-cron-old"),

		testutil.Mock_Events_PutRule_Generic(),
		testutil.Mock_Events_PutTargets_Generic(),
	}

	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd1"), taskmock.WithExitCode(1))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd2"), taskmock.WithExitCode(0))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-sc-inherit"), taskmock.WithExitCode(0))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-storage"))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-override-defaults"))...)
	// mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-disabled"))...)

	for _, familyName := range []string{
		"dummy-pd1", "dummy-pd2", "dummy-pd-sc-inherit", "dummy-pd-storage", "dummy-pd-override-defaults",
		"dummy-console", "dummy-svc1", "dummy-svc2", "dummy-svc3", "dummy-svc4",
		"dummy-cron1", "dummy-cron2", "dummy-cron-daily",
	} {

		revs := []int{998, 999}
		if familyName == "dummy-console" {
			revs = []int{997, 998, 999}
		}

		mocks = append(mocks,
			testutil.Mock_ECS_ListTaskDefinitions(familyName, revs),
			testutil.Mock_ECS_DeregisterTaskDefinition(familyName, 997),
			testutil.Mock_ECS_DeregisterTaskDefinition(familyName, 998),
		)
	}

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: mocks,
	})
}
