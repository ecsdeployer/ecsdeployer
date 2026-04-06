package cmdutil

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/log"
)

func DeprecateWarn(ctx *config.Context) {
	if ctx.Deprecated {
		log.Warn(BoldStyle.Render("you are using deprecated features, check the log above for information"))
	}
}
