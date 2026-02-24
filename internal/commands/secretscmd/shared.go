package secretscmd

import (
	"context"
	"fmt"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

const (
	outputFormatDotEnv = `dotenv`
	outputFormatPlain  = `plain`
)

func loadProject(ctx context.Context, cfgFile string) (*config.Context, string, error) {
	proj, err := cmdutil.LoadConfig(cfgFile)
	if err != nil {
		return nil, "", err
	}

	if !proj.Settings.SSMImport.IsEnabled() {
		return nil, "", fmt.Errorf(`SSM import is not enabled for this project, nothing to list.`)
	}

	ssmImport := *proj.Settings.SSMImport

	cfgCtx := config.Wrap(ctx, proj)

	ssmPrefix, err := tmpl.New(cfgCtx).Apply(ssmImport.GetPath())
	if err != nil {
		return nil, "", err
	}

	// Trim any trailing slash, then add our own
	ssmPrefix = strings.TrimSuffix(ssmPrefix, "/") + "/"

	return cfgCtx, ssmPrefix, nil
}
