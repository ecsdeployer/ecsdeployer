package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// apply all the common stuff
func (b *Builder) applyContainerHealthCheck(cdef *ecsTypes.ContainerDefinition, check *config.HealthCheck) error {

	if check == nil {
		return nil
	}

	if check.Disabled {
		return nil
	}

	value := &ecsTypes.HealthCheck{
		Command: check.Command,
	}

	if check.Interval != nil {
		value.Interval = aws.Int32(check.Interval.ToAwsInt32())
	}
	if check.Retries != nil {
		value.Retries = check.Retries
	}
	if check.StartPeriod != nil {
		value.StartPeriod = aws.Int32(check.StartPeriod.ToAwsInt32())
	}
	if check.Timeout != nil {
		value.Timeout = aws.Int32(check.Timeout.ToAwsInt32())
	}

	cdef.HealthCheck = value

	return nil
}
