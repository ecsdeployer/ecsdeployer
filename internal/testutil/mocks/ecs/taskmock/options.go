package taskmock

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Options struct {
	PendingCount int
	RunningCount int
	ExitCode     int
	Family       string
	StopReason   ecsTypes.TaskStopCode
}

type optFunc = func(*Options)

func WithCounts(pcount, rcount int) optFunc {
	return func(o *Options) {
		o.PendingCount = pcount
		o.RunningCount = rcount
	}
}

func WithFamily(val string) optFunc {
	return func(o *Options) {
		o.Family = val
	}
}

func WithExitCode(val int) optFunc {
	return func(o *Options) {
		o.ExitCode = val
	}
}

func WithSReason(val ecsTypes.TaskStopCode) optFunc {
	return func(o *Options) {
		o.StopReason = val
	}
}
