package predeploytask

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	runTaskBuilder "ecsdeployer.com/ecsdeployer/internal/builders/runtask"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/substep/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/caarlos0/log"
)

type Step struct {
	pdTask *config.PreDeployTask
}

const (
	fieldRuntime = "runtime"
	fieldStatus  = "status"
	fieldReason  = "reason"
	fieldDetail  = "detail"
	fieldTaskArn = "taskArn"
)

var (
	errTaskStillRunning = errors.New("Task is still running")
)

func New(task *config.PreDeployTask) *Step {
	return &Step{
		pdTask: task,
	}
}

func (s *Step) String() string {
	return fmt.Sprintf("task:%s", s.pdTask.Name)
}

func (s *Step) Run(ctx *config.Context) error {

	if s.pdTask.Disabled {
		return step.Skip("task disabled")
	}

	taskDefArn, err := taskdefinition.New(s.pdTask).Register(ctx)
	if err != nil {
		return err
	}

	// logger := log.WithField("name", s.pdTask.Name)
	log.Info("running")

	runTaskInput, err := runTaskBuilder.Build(ctx, s.pdTask)
	if err != nil {
		return err
	}
	runTaskInput.TaskDefinition = aws.String(taskDefArn)

	ecsClient := awsclients.ECSClient()

	startTime := time.Now()

	runTaskOutput, err := ecsClient.RunTask(ctx.Context, runTaskInput)
	if err != nil {
		return err
	}

	if len(runTaskOutput.Failures) > 0 {
		for _, failure := range runTaskOutput.Failures {
			log.WithFields(log.Fields{
				fieldReason: *failure.Reason,
				fieldDetail: *failure.Detail,
			}).Error("failed to launch")
		}
		return errors.New("task failed to launch")
	}

	taskArn := *runTaskOutput.Tasks[0].TaskArn
	log.WithField(fieldTaskArn, taskArn).Info("launched")

	if s.pdTask.DoNotWait {
		log.Warn("skip waiting")
		return nil
	}

	log.Debug("waiting")

	params := &ecs.DescribeTasksInput{
		Tasks:   []string{taskArn},
		Cluster: runTaskInput.Cluster,
	}

	// determine the max wait time either specifically on this task or use the default
	maxWaitTime := util.Coalesce(s.pdTask.Timeout, ctx.Project.Settings.PreDeployTimeout).ToDuration()

	logger := log.WithField("timeout", maxWaitTime)

	// runningWaiter := ecs.NewTasksRunningWaiter(ecsClient, func(trwo *ecs.TasksRunningWaiterOptions) {
	// 	trwo.MinDelay, trwo.MaxDelay = helpers.GetAwsWaiterDelays(5*time.Second, 60*time.Second)
	// })

	// just a dumb wait to make sure the task shows up on AWS API
	if !helpers.IsTestingMode {
		time.Sleep(5 * time.Second)
	}

	stoppedWaiter := ecs.NewTasksStoppedWaiter(ecsClient, func(trwo *ecs.TasksStoppedWaiterOptions) {
		// trwo.MinDelay = 5 * time.Second
		// trwo.MaxDelay = 60 * time.Second

		trwo.MinDelay, trwo.MaxDelay = helpers.GetAwsWaiterDelays(5*time.Second, 60*time.Second)

		oldRetryable := trwo.Retryable
		trwo.Retryable = func(ctx context.Context, dti *ecs.DescribeTasksInput, dto *ecs.DescribeTasksOutput, err error) (bool, error) {

			if err != nil {
				return oldRetryable(ctx, dti, dto, err)
			}

			logger.WithFields(log.Fields{
				fieldRuntime: time.Since(startTime).Round(time.Second).String(),
				fieldStatus:  aws.ToString(dto.Tasks[0].LastStatus),
			}).Trace("waiting...")

			return oldRetryable(ctx, dti, dto, err)
		}

	})

	// wait for task to complete
	err = stoppedWaiter.Wait(ctx.Context, params, maxWaitTime)
	if err != nil {
		log.Error("failed")
		return fmt.Errorf("Failure waiting for task to stop: %w", err)
	}

	// it's stopped, so get the latest status
	results, err := ecsClient.DescribeTasks(ctx.Context, params)
	if err != nil {
		log.Error("failed")
		return fmt.Errorf("Unable to describe task status: %w", err)
	}

	// check for failures
	if len(results.Failures) > 0 {
		for _, failure := range results.Failures {
			log.WithFields(log.Fields{
				fieldReason: aws.ToString(failure.Reason),
				fieldDetail: aws.ToString(failure.Detail),
			}).Error("describe failed")
		}

		if !s.pdTask.IgnoreFailure {
			return errors.New("Task failed to describe")
		}
	}

	result := results.Tasks[0]

	// ensure that there were no task failures (like exit codes or failure to launch)
	err = didTaskSucceed(&result)

	if err == nil {
		log.Info("finished")
		return nil
	}

	if s.pdTask.IgnoreFailure {
		log.WithError(err).Warn("failure (ignored)")
		return nil
	}

	log.WithError(err).Error("failed!")

	return fmt.Errorf("Task failed: %w", err)
}

func didTaskSucceed(result *ecsTypes.Task) error {

	if aws.ToString(result.LastStatus) != string(ecsTypes.DesiredStatusStopped) {
		fmt.Println("TASK: ", util.Must(util.JsonifyPretty(*result)))
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
