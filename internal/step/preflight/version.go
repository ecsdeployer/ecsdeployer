package preflight

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"ecsdeployer.com/ecsdeployer/pkg/version"
	hcVersion "github.com/hashicorp/go-version"
)

type checkVersion struct{}

func (checkVersion) String() string {
	return "version requirements"
}

func (checkVersion) Skip(ctx *config.Context) bool {
	return false
}

func (checkVersion) Check(ctx *config.Context) error {

	reqVersion := ctx.Project.EcsDeployerOptions.RequiredVersion

	if reqVersion == nil {
		return nil
	}

	if isVersionAllowed(reqVersion, version.SemVer) {
		return nil
	}

	return errors.New("Your configuration file prevents this version of ECSDeployer from being used.")
}

func isVersionAllowed(constraints *config.VersionConstraint, currentVersion *hcVersion.Version) bool {
	return constraints.Check(currentVersion)
}
