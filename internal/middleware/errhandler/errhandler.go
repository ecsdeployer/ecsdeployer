package errhandler

import (
	"ecsdeployer.com/ecsdeployer/internal/middleware"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

func Handle(action middleware.Action) middleware.Action {
	return func(ctx *config.Context) error {
		err := action(ctx)
		if err == nil {
			return nil
		}
		if step.IsSkip(err) {
			log.WithField("reason", err.Error()).Warn("step skipped")
			return nil
		}
		return err
	}
}

func Ignore(action middleware.Action) middleware.Action {
	return func(ctx *config.Context) error {
		err := action(ctx)
		if err == nil {
			return nil
		}
		if step.IsSkip(err) {
			log.WithField("reason", err.Error()).Warn("step skipped")
			return nil
		}
		log.WithField("reason", err.Error()).Error("step failed (continuing anyway)")
		return nil
	}
}
