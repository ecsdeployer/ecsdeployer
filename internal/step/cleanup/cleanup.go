package cleanup

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanupcronjobs"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanupservices"
	"ecsdeployer.com/ecsdeployer/internal/step/cleanuptaskdefinitions"
	"ecsdeployer.com/ecsdeployer/internal/step/deregistertaskdefinitions"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type cleaner interface {
	fmt.Stringer
	Clean(ctx *config.Context) error
}

var cleaners = []cleaner{
	deregistertaskdefinitions.Step{},
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
	for _, cleaner := range cleaners {
		if err := skip.Maybe(
			cleaner,
			logging.PadLog(
				cleaner.String(),
				errhandler.Ignore(cleaner.Clean),
			),
		)(ctx); err != nil {

			return fmt.Errorf("%s: failed with: %w", cleaner.String(), err)
		}
	}
	return nil
}
