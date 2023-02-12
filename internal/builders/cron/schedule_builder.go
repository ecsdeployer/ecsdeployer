package cron

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	scheduler "github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func BuildSchedule(ctx *config.Context, resource *config.CronJob, taskDefArn string) (*scheduler.CreateScheduleInput, error) {

	project := ctx.Project
	templates := project.Templates

	// This does not reference job specific values, so does not use the same template
	scheduleGroupName, err := tmpl.New(ctx).Apply(*templates.ScheduleGroupName)
	if err != nil {
		return nil, err
	}

	tpl := tmpl.New(ctx).WithExtraFields(resource.TemplateFields())

	scheduleName, err := tpl.Apply(*templates.ScheduleName)
	if err != nil {
		return nil, err
	}

	cronContainerName, err := tpl.Apply(*templates.ContainerName)
	if err != nil {
		return nil, err
	}

	ecsParams := &schedulerTypes.EcsParameters{
		TaskDefinitionArn:    aws.String(taskDefArn),
		TaskCount:            aws.Int32(1),
		EnableECSManagedTags: aws.Bool(true),
		EnableExecuteCommand: aws.Bool(false),
		LaunchType:           schedulerTypes.LaunchTypeFargate,
		PlatformVersion:      resource.PlatformVersion,
		PropagateTags:        schedulerTypes.PropagateTagsTaskDefinition,

		// this isnt needed since the launchtype is fargate
		// CapacityProviderStrategy: config.NewSpotOnDemand().ExportCapacityStrategyScheduler(),
	}

	cronGroupName, err := tpl.Apply(*templates.CronGroup)
	if err != nil {
		return nil, err
	}
	if cronGroupName != "" {
		ecsParams.Group = aws.String(cronGroupName)
	}

	// Cronjob Input field
	cronInput := cronInputObj{}

	if !project.Settings.SkipCronEnvVars {
		cronEnvVars := make([]cronOverrideKeyPair, 0, len(config.DefaultCronEnvVars))
		for k, v := range config.DefaultCronEnvVars {
			cronEnvVars = append(cronEnvVars, cronOverrideKeyPair{
				Name:  k,
				Value: v,
			})
		}
		cronInput.ContainerOverrides = []cronContainerOverride{
			{
				Name:        cronContainerName,
				Environment: cronEnvVars,
			},
		}
	}

	cronInputJson, err := util.Jsonify(cronInput)
	if err != nil {
		return nil, err
	}

	clusterArn, err := project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	// Load network configuration
	network := util.Coalesce(resource.Network, project.TaskDefaults.Network, project.Network)
	if network == nil {
		return nil, errors.New("Unable to resolve network configuration!")
	}

	ecsNetworkConfig, err := network.ResolveSched(ctx)
	if err != nil {
		return nil, err
	}
	ecsParams.NetworkConfiguration = ecsNetworkConfig

	// The target
	targetParams := &schedulerTypes.Target{
		Arn:           aws.String(clusterArn),
		Input:         aws.String(cronInputJson),
		EcsParameters: ecsParams,
		RetryPolicy: &schedulerTypes.RetryPolicy{
			MaximumEventAgeInSeconds: aws.Int32(180),
			MaximumRetryAttempts:     aws.Int32(0),
		},
	}

	if project.CronLauncherRole != nil {
		launcherRole, err := project.CronLauncherRole.Arn(ctx)
		if err != nil {
			return nil, err
		}
		targetParams.RoleArn = &launcherRole
	}

	flexWindow := &schedulerTypes.FlexibleTimeWindow{
		Mode: schedulerTypes.FlexibleTimeWindowModeOff,
	}

	scheduleInput := &scheduler.CreateScheduleInput{
		FlexibleTimeWindow: flexWindow,
		State:              schedulerTypes.ScheduleStateEnabled,
		Name:               aws.String(scheduleName),
		GroupName:          aws.String(scheduleGroupName),
		Target:             targetParams,
		Description:        aws.String(resource.Description),
		ScheduleExpression: aws.String(resource.Schedule),
		StartDate:          resource.StartDate,
		EndDate:            resource.EndDate,
	}

	if resource.TimeZone != nil {
		scheduleInput.ScheduleExpressionTimezone = resource.TimeZone
	}

	if resource.IsDisabled() {
		scheduleInput.State = schedulerTypes.ScheduleStateDisabled
	}

	return scheduleInput, nil
}
