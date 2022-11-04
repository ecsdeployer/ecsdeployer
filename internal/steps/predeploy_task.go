package steps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/builders"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/caarlos0/log"
)

var (
	errTaskStillRunning = errors.New("Task is still running")
)

const (
	fieldRuntime = "runtime"
	fieldStatus  = "status"
	fieldReason  = "reason"
	fieldDetail  = "detail"
	fieldTaskArn = "taskArn"
)

func PreDeployTaskStep(resource *config.PreDeployTask) *Step {

	if resource.Disabled {
		// task was disabled, bail out
		return NoopStep()
	}

	return NewStep(&Step{
		Label:    "PreDeployTask",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepPreDeployTaskCreate,
		Dependencies: []*Step{
			TaskDefinitionStep(resource),
		},
	})
}

func stepPreDeployTaskCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {
	taskDefOutput, ok := step.LookupOutput("task_definition_arn")
	if !ok {
		return nil, errors.New("task definition was not created, cannot run predeploy")
	}
	taskDefArn := (taskDefOutput).(string)

	predeployTask := (step.Resource).(*config.PreDeployTask)
	logger := step.Logger

	runTaskInput, err := builders.BuildRunTask(ctx, predeployTask)
	if err != nil {
		return nil, err
	}
	runTaskInput.TaskDefinition = aws.String(taskDefArn)

	logger.Info("Running PreDeployTask")

	ecsClient := ctx.ECSClient()

	startTime := time.Now()

	runTaskOutput, err := ecsClient.RunTask(ctx.Context, runTaskInput)
	if err != nil {
		return nil, err
	}

	if len(runTaskOutput.Failures) > 0 {
		for _, failure := range runTaskOutput.Failures {
			logger.WithFields(log.Fields{
				fieldReason: aws.ToString(failure.Reason),
				fieldDetail: aws.ToString(failure.Detail),
			}).Error("Task Failed to Launch")
		}
		return nil, errors.New("task failed to launch")
	}

	taskArn := aws.ToString(runTaskOutput.Tasks[0].TaskArn)
	logger = logger.WithField(fieldTaskArn, taskArn)
	logger.Info("Task launched")

	fields := OutputFields{
		fieldTaskArn: taskArn,
	}

	logger.Debugf("Waiting for task to complete")

	params := &ecs.DescribeTasksInput{
		Tasks:   []string{taskArn},
		Cluster: runTaskInput.Cluster,
	}

	// determine the max wait time either specifically on this task or use the default
	maxWaitTime := util.Coalesce(predeployTask.Timeout, ctx.Project.Settings.PreDeployTimeout).ToDuration()

	logger = logger.WithField("timeout", maxWaitTime)

	// runningWaiter := ecs.NewTasksRunningWaiter(ecsClient, func(trwo *ecs.TasksRunningWaiterOptions) {
	// 	trwo.MinDelay = 10 * time.Second
	// 	trwo.MaxDelay = 30 * time.Second
	// })

	// just a dumb wait to make sure the task shows up on AWS API
	time.Sleep(5 * time.Second)

	stoppedWaiter := ecs.NewTasksStoppedWaiter(ecsClient, func(trwo *ecs.TasksStoppedWaiterOptions) {
		trwo.MinDelay = 10 * time.Second
		trwo.MaxDelay = 60 * time.Second

		oldRetryable := trwo.Retryable
		trwo.Retryable = func(ctx context.Context, dti *ecs.DescribeTasksInput, dto *ecs.DescribeTasksOutput, err error) (bool, error) {
			logger.WithFields(log.Fields{
				fieldRuntime: time.Since(startTime).Round(time.Second).String(),
				fieldStatus:  aws.ToString(dto.Tasks[0].LastStatus),
			}).Info("Waiting for task...")

			return oldRetryable(ctx, dti, dto, err)
		}

	})

	// wait for task to complete
	err = stoppedWaiter.Wait(ctx.Context, params, maxWaitTime)
	if err != nil {
		logger.Error("Failed")
		return nil, err
	}

	// it's stopped, so get the latest status
	results, err := ecsClient.DescribeTasks(ctx.Context, params)
	if err != nil {
		logger.Error("Failed")
		return nil, err
	}

	// check for failures
	if len(results.Failures) > 0 {
		for _, failure := range results.Failures {
			logger.WithFields(log.Fields{
				fieldReason: aws.ToString(failure.Reason),
				fieldDetail: aws.ToString(failure.Detail),
			}).Error("Task describe failed")
		}

		if !predeployTask.IgnoreFailure {
			return nil, errors.New("Task failed to describe")
		}
	}

	result := results.Tasks[0]

	// ensure that there were no task failures (like exit codes or failure to launch)
	err = didTaskSucceed(&result)

	if err == nil {
		logger.Info("Task completed successfully!")
		return fields, nil
	}

	if predeployTask.IgnoreFailure {
		logger.WithError(err).Warn("Task Failed, but failures are ignored")
		return fields, nil
	}

	logger.Error("Task Failed")

	return fields, err
}

func didTaskSucceed(result *ecsTypes.Task) error {

	if aws.ToString(result.LastStatus) != string(ecsTypes.DesiredStatusStopped) {
		return errTaskStillRunning
	}

	stopCode := result.StopCode

	if stopCode == ecsTypes.TaskStopCodeTaskFailedToStart {
		return fmt.Errorf("Failed to Start: %s", aws.ToString(result.StoppedReason))
	}

	if stopCode == ecsTypes.TaskStopCodeUserInitiated {
		return errors.New("User killed the task")
	}

	if stopCode != ecsTypes.TaskStopCodeEssentialContainerExited {
		return fmt.Errorf("Some very weird stop code was given: %s", string(stopCode))
	}

	for _, cont := range result.Containers {
		if cont.ExitCode != nil {
			exitCode := aws.ToInt32(cont.ExitCode)
			if exitCode > 0 {
				return fmt.Errorf("Container exited with code: %d", exitCode)
			}
		}
	}

	return nil
}
