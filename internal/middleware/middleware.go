package middleware

import "ecsdeployer.com/ecsdeployer/pkg/config"

type Action func(ctx *config.Context) error
