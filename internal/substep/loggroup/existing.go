package loggroup

import (
	"ecsdeployer.com/ecsdeployer/internal/step/preloadloggroups"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// 'forced' - bypass any cache check and request from AWS
func (s *Substep) describeLogGroup(ctx *config.Context, forced bool) (*logTypes.LogGroup, error) {

	if forced {
		err := preloadloggroups.ByPrefix(ctx, s.groupName)
		if err != nil {
			return nil, err
		}
	}

	if val, ok := ctx.Cache.LogGroups[s.groupName]; ok {
		return &val, nil
	}

	return nil, nil
}
