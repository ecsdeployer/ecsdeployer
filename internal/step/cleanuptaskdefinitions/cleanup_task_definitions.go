// this will deregister task definitions that are no longer being managed by ECSDeployer
package cleanuptaskdefinitions

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned task definitions"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsEnabled
}

func (Step) Clean(ctx *config.Context) error {

	return nil
}
