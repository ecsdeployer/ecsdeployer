package cmd

import (
	"context"
	"errors"
	"os"
	"time"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

var genericConfigPaths = [...]string{
	".ecsdeployer.yml",
	".ecsdeployer.yaml",
	// "ecsdeployer.yml",
	// "ecsdeployer.yaml",
}

type configLoaderExtras struct {
	configFile  string
	appVersion  string
	imageTag    string
	imageUri    string
	noValidate  bool
	timeout     time.Duration
	cmdMetadata *cmdMetadata
}

func loadProjectContext(options *configLoaderExtras) (*config.Context, context.CancelFunc, error) {
	cfg, err := loadConfig(options.configFile)
	if err != nil {
		return nil, nil, err
	}

	if options.cmdMetadata != nil {
		if err = checkVersionRestriction(cfg, options.cmdMetadata); err != nil {
			return nil, nil, err
		}
	}

	if options.timeout == 0 {
		options.timeout = 90 * time.Minute
	}

	ctx, cancel := config.NewWithTimeout(cfg, options.timeout)
	// defer cancel()

	ctx.Version = options.appVersion
	ctx.ImageTag = options.imageTag

	if cfg.StageName != nil {
		ctx.Stage = *cfg.StageName
	}

	// override the image
	if options.imageUri != "" {
		ctx.ImageUriRef = options.imageUri
	}

	if err := ctx.Project.ValidateWithContext(ctx); err != nil {
		defer cancel()
		return nil, nil, err
	}
	return ctx, cancel, nil
}

func loadConfig(path string) (*config.Project, error) {
	if path == "-" {
		log.Info("loading config from stdin")
		return config.LoadReader(os.Stdin)
	}
	if path != "" {
		return config.Load(path)
	}

	for _, filepath := range genericConfigPaths {
		if _, err := os.Stat(filepath); err == nil {
			log.WithField("config", filepath).Debug("Using default configuration file")
			return config.Load(filepath)
		}
	}

	return nil, errors.New("No configuration file was found")
}

func checkVersionRestriction(config *config.Project, metadata *cmdMetadata) error {
	ok, errorList := config.EcsDeployerOptions.IsVersionAllowed(metadata.version)
	if len(errorList) > 0 {
		for _, err := range errorList {
			log.WithError(err).Error("Version Restriction Failed")
		}
		return errors.New("Version restriction failed")
	}

	if !ok {
		return errors.New("Your configuration file prevents this version of ECSDeployer from being used.")
	}

	return nil
}
