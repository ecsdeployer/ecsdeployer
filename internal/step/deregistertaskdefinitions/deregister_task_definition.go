// This deregisters task definitions that are still being managed by ecsdeployer
package deregistertaskdefinitions

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

type Step struct{}

func (Step) String() string {
	return "cleaning previous task definitions"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsEnabled &&
		ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsOnlyManaged
}

func (Step) Clean(ctx *config.Context) error {

	for _, defArn := range ctx.Cache.RegisteredTaskDefArns {
		family := helpers.GetTaskDefFamilyFromArn(defArn)
		log.WithField("family", family).Debug(defArn)
	}

	return nil
}
