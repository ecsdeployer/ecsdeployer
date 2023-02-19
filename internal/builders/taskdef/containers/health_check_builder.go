package containers

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type HealthCheckBuilder struct {
	check *config.HealthCheck
}

func (hc *HealthCheckBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {

	if hc.check == nil {
		return nil
	}

	value := &ecsTypes.HealthCheck{
		Command: hc.check.Command,
	}

	if hc.check.Interval != nil {
		value.Interval = aws.Int32(hc.check.Interval.ToAwsInt32())
	}
	if hc.check.Retries != nil {
		value.Retries = hc.check.Retries
	}
	if hc.check.StartPeriod != nil {
		value.StartPeriod = aws.Int32(hc.check.StartPeriod.ToAwsInt32())
	}
	if hc.check.Timeout != nil {
		value.Timeout = aws.Int32(hc.check.Timeout.ToAwsInt32())
	}

	obj.Data.HealthCheck = value

	return nil
}
