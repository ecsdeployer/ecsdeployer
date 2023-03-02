package steps

import (
	"context"
	"errors"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	serviceBuilder "ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const (
	serviceStepAttrName = "ServiceName"
)

var (
	ErrTaskDefNotCreated = errors.New("Task definition was not created")
)

func ServiceStep(resource *config.Service) *Step {
	return NewStep(&Step{
		Label:    "Service",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepServiceCreate,
		Read:     stepServiceRead,
		Update:   stepServiceUpdate,
		PreApply: stepServicePreApply,
		Dependencies: []*Step{
			TaskDefinitionStep(resource),
		},
	})
}

func stepServicePreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {

	common, err := config.ExtractCommonTaskAttrs(step.Resource)
	if err != nil {
		return err
	}

	tplFields, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	if err != nil {
		return err
	}

	serviceName, err := tmpl.New(ctx).WithExtraFields(tplFields).Apply(*ctx.Project.Templates.ServiceName)
	if err != nil {
		return err
	}

	step.SetAttr(serviceStepAttrName, serviceName)

	return nil
}

func stepServiceCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {
	logger := step.Logger
	logger.Info("Creating service")

	svc := (step.Resource).(*config.Service)

	taskDefOutput, ok := step.LookupOutput("task_definition_arn")
	if !ok {
		return nil, ErrTaskDefNotCreated
	}

	createServiceInput, err := serviceBuilder.BuildCreate(ctx, svc)
	if err != nil {
		return nil, err
	}

	createServiceInput.TaskDefinition = aws.String(taskDefOutput.(string))

	result, err := awsclients.ECSClient().CreateService(ctx.Context, createServiceInput)
	if err != nil {
		return nil, err
	}

	fields := OutputFields{
		"Service": result.Service,
	}

	err = stepServiceWaitForSuccess(ctx, step, result.Service)
	if err != nil {
		return fields, err
	}

	return fields, nil
}

func stepServiceUpdate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {
	step.Logger.Info("Updating service")

	svc := (step.Resource).(*config.Service)

	taskDefOutput, ok := step.LookupOutput("task_definition_arn")
	if !ok {
		return nil, ErrTaskDefNotCreated
	}

	updateServiceInput, err := serviceBuilder.BuildUpdate(ctx, svc)
	if err != nil {
		return nil, err
	}

	updateServiceInput.TaskDefinition = aws.String(taskDefOutput.(string))

	result, err := awsclients.ECSClient().UpdateService(ctx.Context, updateServiceInput)
	if err != nil {
		return nil, err
	}

	fields := OutputFields{
		"Service": result.Service,
	}

	err = stepServiceWaitForSuccess(ctx, step, result.Service)
	if err != nil {
		return fields, err
	}

	return fields, nil
}

func stepServiceRead(ctx *config.Context, step *Step, meta *StepMetadata) (any, error) {

	serviceNameRes, ok := step.GetAttr(serviceStepAttrName)
	if !ok {
		// the attribute key is missing, that means this should be skipped
		return nil, errors.New("Service Name Key is missing")
	}
	serviceName := serviceNameRes.(string)

	clusterArn, err := ctx.Project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	ecsClient := awsclients.ECSClient()
	result, err := ecsClient.DescribeServices(ctx.Context, &ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Cluster:  aws.String(clusterArn),
	})
	if err != nil {
		return nil, err
	}

	if len(result.Failures) > 0 {
		failReason := aws.ToString(result.Failures[0].Reason)
		if failReason == "MISSING" {
			return nil, nil
		}

		return nil, errors.New(failReason)
	}

	svc := result.Services[0]

	return svc, nil
}

func stepServiceWaitForSuccess(ctx *config.Context, step *Step, service *ecsTypes.Service) error {

	waitForStable := ctx.Project.Settings.WaitForStable

	logger := step.Logger

	if !*waitForStable.Individually {
		// let it be handled by the ServiceDeployment step
		return nil
	}

	if waitForStable.IsDisabled() {
		step.Logger.Warn("You have requested to skip stability checks.")
		return nil
	}

	ecsClient := awsclients.ECSClient()
	startTime := time.Now()

	waiter := ecs.NewServicesStableWaiter(ecsClient, func(sswo *ecs.ServicesStableWaiterOptions) {
		sswo.MinDelay, sswo.MaxDelay = helpers.GetAwsWaiterDelays(10*time.Second, 45*time.Second)
		sswo.LogWaitAttempts = false

		oldRetryable := sswo.Retryable
		sswo.Retryable = func(ctx context.Context, dsi *ecs.DescribeServicesInput, dso *ecs.DescribeServicesOutput, err error) (bool, error) {

			if err != nil {
				return false, err
			}

			logger.WithField(fieldRuntime, time.Since(startTime).Round(time.Second).String()).Info("Waiting for service...")

			return oldRetryable(ctx, dsi, dso, err)
		}
	})

	params := &ecs.DescribeServicesInput{
		Services: []string{*service.ServiceName},
		Cluster:  service.ClusterArn,
	}

	maxWaitTime := ctx.Project.Settings.WaitForStable.Timeout.ToDuration()

	err := waiter.Wait(ctx.Context, params, maxWaitTime)
	if err != nil {
		logger.Error("Service unstable!")
		return err
	}

	logger.Info("Service stable")

	return nil
}
