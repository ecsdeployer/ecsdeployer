package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSmokeTest_StepFuncs(t *testing.T) {
	project := config.Project{
		ConsoleTask: &config.ConsoleTask{
			Enabled: util.Ptr(true),
		},
		PreDeployTasks: []*config.PreDeployTask{
			{
				CommonTaskAttrs: config.CommonTaskAttrs{
					CommonContainerAttrs: config.CommonContainerAttrs{
						Name: "pd1",
					},
				},
			},
		},
		CronJobs: []*config.CronJob{
			{
				CommonTaskAttrs: config.CommonTaskAttrs{
					CommonContainerAttrs: config.CommonContainerAttrs{
						Name: "cron1",
					},
				},
			},
		},
		Services: []*config.Service{
			{
				CommonTaskAttrs: config.CommonTaskAttrs{
					CommonContainerAttrs: config.CommonContainerAttrs{
						Name: "svc1",
					},
				},
			},
		},
	}
	project.ApplyDefaults()

	tables := []struct {
		step  *Step
		label string
	}{
		{CleanupStep(project.Settings.KeepInSync), "Cleanup"},
		{ConsoleTaskStep(project.ConsoleTask), "ConsoleTask"},
		{CleanupTaskDefinitionsStep(project.Settings.KeepInSync), "CleanupTaskDefinitions"},
		{CleanupCronjobsStep(project.Settings.KeepInSync), "CleanupCronjobs"},
		{CleanupServicesStep(project.Settings.KeepInSync), "CleanupServices"},
		{CleanupOnlyStep(&project), "CleanupOnly"},
		{CronDeploymentStep(&project), "CronDeployment"},
		// {CronRuleStep(project.CronJobs[0]), "CronRule"},
		// {CronTargetStep(project.CronJobs[0]), "CronTarget"},
		{CronSchedulesStep(&project), "CronSchedules"},
		{CronjobStep(project.CronJobs[0]), "Cronjob"},
		{DeploymentStep(&project), "Deployment"},
		{DeregisterTaskDefinitionsStep(&project), "DeregisterTaskDefinitions"},
		{FirelensStep(&project), "Firelens"},
		{LogGroupStep(project.ConsoleTask), "LogGroup"},
		{NoopStep(), "Noop"},
		{PreDeployTaskStep(project.PreDeployTasks[0]), "PreDeployTask"},
		{PreDeploymentStep(&project), "PreDeployment"},
		{PreflightStep(&project), "Preflight"},
		{PreloadLogGroupsStep(&project), "PreloadLogGroups"},
		{PreloadSecretsStep(&project), "PreloadSecrets"},
		{PreloadStep(&project), "Preload"},
		{ScheduleGroupStep(&project), "ScheduleGroup"},
		{ServiceDeploymentStep(&project), "ServiceDeployment"},
		{ServiceStep(project.Services[0]), "Service"},
		{TargetGroupStep(project.Services[0]), "TargetGroup"},
		{TaskDefinitionStep(project.ConsoleTask), "TaskDefinition"},
	}

	for _, table := range tables {
		require.Equal(t, table.label, table.step.Label)
	}

}
