package servicedeployment

import (
	"context"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/caarlos0/log"
)

func waitForStable(ctx *config.Context, service *ecsTypes.Service) error {

	waitForStable := ctx.Project.Settings.WaitForStable

	// if !*waitForStable.Individually {
	// 	// let it be handled by the ServiceDeployment step
	// 	return nil
	// }

	logger := log.WithField("name", *service.ServiceName)

	if waitForStable.IsDisabled() {
		logger.Info("skipping stability checks")
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

			logger.WithField("runtime", time.Since(startTime).Round(time.Second).String()).Trace("waiting for stable")

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
		logger.Error("service unstable")
		return err
	}

	logger.Info("deployed")

	return nil
}
