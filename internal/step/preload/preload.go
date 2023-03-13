package preload

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/step/preloadloggroups"
	"ecsdeployer.com/ecsdeployer/internal/step/preloadsecrets"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type preloader interface {
	fmt.Stringer
	Preload(ctx *config.Context) error
}

var preloaders = []preloader{
	preloadsecrets.Step{},
	preloadloggroups.Step{},
}

type Step struct{}

func (Step) String() string {
	return "preloading resources"
}

func (Step) Skip(ctx *config.Context) bool { return false }

func (Step) Run(ctx *config.Context) error {
	for _, preloader := range preloaders {
		if err := skip.Maybe(
			preloader,
			logging.PadLog(
				preloader.String(),
				errhandler.Handle(preloader.Preload),
			),
		)(ctx); err != nil {
			return fmt.Errorf("%s: failed with: %w", preloader.String(), err)
		}
	}
	return nil
}
