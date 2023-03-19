package cmd

import (
	"fmt"
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/steps"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	dsmock "ecsdeployer.com/ecsdeployer/internal/testutil/mocks/ecs/describeservicemock"
	"ecsdeployer.com/ecsdeployer/internal/testutil/mocks/ecs/taskmock"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestDeployCmd(t *testing.T) {
	silenceLogging(t)

	t.Run("calls correct function", func(t *testing.T) {
		oldRef := stepDeploymentStepFunc
		t.Cleanup(func() {
			stepDeploymentStepFunc = oldRef
		})

		testutil.StartMocker(t, nil)

		wasCalled := false
		stepDeploymentStepFunc = func(_ *config.Project) *steps.Step {
			wasCalled = true
			return steps.NoopStep()
		}

		cmd := newDeployCmd(defaultCmdMetadata()).cmd
		cmd.Root().SetArgs([]string{"-q"})
		cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

		_, _, err := executeCmdAndReturnOutput(cmd)

		require.NoError(t, err)

		require.True(t, wasCalled)
	})
}

func TestDeploySmoke(t *testing.T) {
	helpers.IsTestingMode = true

	orig := log.Log
	t.Cleanup(func() {
		log.Log = orig
	})
	log.Log = log.New(os.Stdout)

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
		testutil.Mock_SSM_GetParametersByPath("/ecsdeployer/dummy/", []string{"SSM_VAR1", "SSM_VAR2"}),
		testutil.Mock_Tagging_GetResources("ecs:task-definition", map[string]string{"ecsdeployer/project": "dummy"}, []string{
			fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-something-something:122", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
			fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-something-something:123", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
		}),
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
	}

	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd1"), taskmock.WithExitCode(1))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd2"), taskmock.WithExitCode(1))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-sc-inherit"), taskmock.WithExitCode(1))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-storage"))...)
	mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-override-defaults"))...)
	// mocks = append(mocks, taskmock.Mock(taskmock.WithFamily("dummy-pd-disabled"))...)

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: mocks,
	})

	lipgloss.SetColorProfile(termenv.TrueColor)
	log.SetLevel(log.TraceLevel)
	cmd := newRootCmd("fake", func(i int) {}).cmd
	log.Strings[log.DebugLevel] = "%"

	// cmd := newDeployCmd(defaultCmdMetadata()).cmd
	// cmd.Root().SetArgs([]string{"-q"})
	// log.SetLevel(log.DebugLevel)
	// cmd.Root().SetArgs([]string{"--debug"})
	// cmd.SetArgs([]string{"deploy", "-c", "../internal/builders/testdata/smoke.yml", "--debug"})
	cmd.SetArgs([]string{"deploy", "-c", "../internal/builders/testdata/everything.yml", "--debug"})

	_, _, err := executeCmdAndReturnOutput(cmd)
	require.NoError(t, err)

}
