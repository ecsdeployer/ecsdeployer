package cmd

import (
	"errors"
	"os"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

var genericConfigPaths = [...]string{
	".ecsdeployer.yml",
	".ecsdeployer.yaml",
	"ecsdeployer.yml",
	"ecsdeployer.yaml",
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
