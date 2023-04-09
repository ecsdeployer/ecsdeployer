package preflight

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/webdestroya/go-log"
)

type Step struct{}

type preflightChecker interface {
	fmt.Stringer
	Check(*config.Context) error
}

var preflightChecks = []preflightChecker{
	checkProject{},         // project validation
	checkVersion{},         // version restrictions
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

		if err := runPFCheck(check, ctx); err != nil {
			return err
		}
	}
	return nil
}

func runPFCheck(check preflightChecker, ctx *config.Context) error {
	if skipcheck, ok := check.(skip.Skipper); ok {
		if skipcheck.Skip(ctx) {
			log.Tracef("skipping %s", check.String())
			return nil
		}
	}

	log.Infof("checking %s", check.String())
	// if err := check.Check(ctx); err != nil {
	// 	return fmt.Errorf("%s: check failed: %w", check.String(), err)
	// }
	if err := skip.Maybe(
		check,
		errhandler.Handle(check.Check),
	)(ctx); err != nil {
		return fmt.Errorf("%s: failed with: %w", check.String(), err)
	}

	return nil
}
