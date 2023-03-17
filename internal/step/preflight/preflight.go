package preflight

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

type Step struct{}

type preflightChecker interface {
	fmt.Stringer
	Check(*config.Context) error
}

var preflightChecks = []preflightChecker{
	checkAccount{},         // acct restriction
	checkTemplates{},       // templates
	checkContainerImages{}, // container images
	checkCluster{},         // cluster
	checkRoles{},           // roles
	checkTargetGroups{},    // target groups
	checkNetwork{},         // network
}

func (Step) String() string {
	return "preflight checks"
}

func (Step) Skip(ctx *config.Context) bool { return false }

func (Step) Run(ctx *config.Context) error {
	for _, check := range preflightChecks {
		log.Infof("checking %s", check.String())
		if err := check.Check(ctx); err != nil {
			return fmt.Errorf("%s: check failed: %w", check.String(), err)
		}
	}
	return nil
}
