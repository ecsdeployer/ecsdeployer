package cmdutil

import (
	"os"

	"ecsdeployer.com/ecsdeployer/internal/usererr"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/webdestroya/go-log"
)

var genericConfigPaths = [...]string{
	".ecsdeployer.yml",
	".ecsdeployer.yaml",
	"ecsdeployer.yml",
	"ecsdeployer.yaml",
}

func LoadConfig(path string) (*config.Project, error) {
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

	return nil, usererr.New("No configuration file was found")
}
