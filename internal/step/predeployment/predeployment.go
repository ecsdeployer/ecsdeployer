package predeployment

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/step/predeploytask"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "predeploy tasks"
}

func (Step) Skip(ctx *config.Context) bool {
	return len(ctx.Project.PreDeployTasks) == 0
}

func (Step) Run(ctx *config.Context) error {

	for _, pdtask := range ctx.Project.PreDeployTasks {
		pdtask := pdtask
		runner := predeploytask.New(pdtask)

		if err := skip.Maybe(
			runner,
			logging.PadLog(
				runner.String(),
				errhandler.Handle(runner.Run),
			),
		)(ctx); err != nil {
			return fmt.Errorf("%s: failed with: %w", runner.String(), err)
		}

		// if err := predeploytask.New(pdtask).Run(ctx); err != nil {
		// 	return err
		// }
	}

	return nil
}
