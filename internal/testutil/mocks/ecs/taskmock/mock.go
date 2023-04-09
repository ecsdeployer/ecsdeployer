package taskmock

import (
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/webdestroya/awsmocker"
)

func Mock(opts ...optFunc) []*awsmocker.MockedEndpoint {

	taskId := testutil.RandomHex(32)

	options := &Options{
		PendingCount: 2,
		RunningCount: 2,
		ExitCode:     0,
		StopReason:   ecsTypes.TaskStopCodeEssentialContainerExited,
	}
	for _, optFunc := range opts {
		optFunc(options)
	}

	lc := &lifecycle{
		options: options,
		taskId:  taskId,
	}

	return []*awsmocker.MockedEndpoint{
		mockRunTask(options.Family, taskId),
		lc.describeTaskMock(),
	}
}
