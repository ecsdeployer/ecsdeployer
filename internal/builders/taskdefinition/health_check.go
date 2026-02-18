package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
		value.Interval = new(check.Interval.ToAwsInt32())
	}
	if check.Retries != nil {
		value.Retries = check.Retries
	}
	if check.StartPeriod != nil {
		value.StartPeriod = new(check.StartPeriod.ToAwsInt32())
	}
	if check.Timeout != nil {
		value.Timeout = new(check.Timeout.ToAwsInt32())
	}

	cdef.HealthCheck = value

	return nil
}
