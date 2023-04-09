package service

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/webdestroya/go-log"
)

func (s *Step) getExisting(ctx *config.Context) (bool, error) {

	log.WithField("name", s.name).Trace("checking existing service")

	clusterArn, err := ctx.Project.Cluster.Arn(ctx)
	if err != nil {
		return false, err
	}

	ecsClient := awsclients.ECSClient()
	result, err := ecsClient.DescribeServices(ctx.Context, &ecs.DescribeServicesInput{
		Services: []string{s.name},
		Cluster:  &clusterArn,
	})
	if err != nil {
		return false, err
	}

	if len(result.Failures) > 0 {
		failReason := *result.Failures[0].Reason
		if failReason == "MISSING" {
			return false, nil
		}

		return false, fmt.Errorf("Unable to describe service: %s", failReason)
	}

	if len(result.Services) == 0 {
		return false, nil
	}

	svc := result.Services[0]

	return *svc.ServiceName == s.name, nil
}
