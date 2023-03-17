package loggroup

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (s *Substep) updateLogGroup(ctx *config.Context, current *logTypes.LogGroup) error {
	return nil
}
