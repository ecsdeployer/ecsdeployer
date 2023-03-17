package servicedeployment

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/caarlos0/log"
)

func describeService(ctx *config.Context, service *config.Service) (*ecsTypes.Service, error) {
	serviceName, err := getServiceName(ctx, service)
	if err != nil {
		return nil, err
	}
	log.WithField("name", serviceName).Trace("checking existing service")

	clusterArn, err := ctx.Project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	ecsClient := awsclients.ECSClient()
	result, err := ecsClient.DescribeServices(ctx.Context, &ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Cluster:  &clusterArn,
	})
	if err != nil {
		return nil, err
	}

	if len(result.Failures) > 0 {
		failReason := *result.Failures[0].Reason
		if failReason == "MISSING" {
			return nil, nil
		}

		return nil, fmt.Errorf("Unable to describe service: %s", failReason)
	}

	return &result.Services[0], nil
}
