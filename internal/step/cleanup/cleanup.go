package cleanup

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanupcronjobs"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanupservices"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanuptaskdefinitions"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type cleaner interface {
	fmt.Stringer
	Clean(ctx *config.Context) error
}

// These cleaner steps must be completely self-sufficient, and must not rely on the cache
var cleaners = []cleaner{
	cleanupservices.Step{},
	cleanuptaskdefinitions.Step{},
	cleanupcronjobs.Step{},
}

type Step struct{}

func (Step) String() string {
	return "cleanup"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Settings.KeepInSync.AllDisabled()
}

func (Step) Run(ctx *config.Context) error {

	wrapRunner := errhandler.Ignore

	if ctx.CleanOnlyFlow {
		wrapRunner = errhandler.Handle
	}

	for _, cleaner := range cleaners {
		if err := skip.Maybe(
			cleaner,
			logging.PadLog(
				cleaner.String(),
				wrapRunner(cleaner.Clean),
			),
		)(ctx); err != nil {

			return fmt.Errorf("%s: failed with: %w", cleaner.String(), err)
		}
	}
	return nil
}
