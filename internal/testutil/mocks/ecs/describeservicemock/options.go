package describeservicemock

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Options struct {
	MaxCount int
	Missing  bool
	Name     string
	Service  ecsTypes.Service
}

type optFunc = func(*Options)

func WithMaxCount(val int) optFunc {
	return func(o *Options) {
		o.MaxCount = val
	}
}

func WithName(val string) optFunc {
	return func(o *Options) {
		o.Name = val
	}
}

func WithMissing() optFunc {
	return func(o *Options) {
		o.Missing = true
	}
}

func WithPending() optFunc {
	return func(o *Options) {
		o.Service.RunningCount = 4
		o.Service.DesiredCount = 2
		o.Service.PendingCount = 1
		o.Service.Deployments = []ecsTypes.Deployment{
			{
				RunningCount: 2,
				DesiredCount: 2,
				PendingCount: 0,
				Status:       aws.String("PRIMARY"),
			},
			{
				RunningCount: 1,
				DesiredCount: 2,
				PendingCount: 1,
				Status:       aws.String("ACTIVE"),
			},
		}
	}
}

func WithStable() optFunc {
	return func(o *Options) {
		o.Service.RunningCount = 2
		o.Service.DesiredCount = 2
		o.Service.PendingCount = 0
		o.Service.Deployments = []ecsTypes.Deployment{
			{
				RunningCount: 2,
				DesiredCount: 2,
				PendingCount: 0,
				Status:       aws.String("PRIMARY"),
			},
		}
	}
}
