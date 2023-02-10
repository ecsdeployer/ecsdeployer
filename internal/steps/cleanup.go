package steps

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type markerTag struct {
	key   string
	value string
}

func CleanupStep(resource *config.KeepInSync) *Step {

	if resource.AllDisabled() {
		return NoopStep()
	}

	return NewStep(&Step{
		Label:        "Cleanup",
		Resource:     resource,
		ParallelDeps: true,
		Dependencies: []*Step{
			CleanupServicesStep(resource),
			CleanupTaskDefinitionsStep(resource),
			CleanupCronjobsStep(resource),
		},
	})
}

func stepCleanupMarkerTag(ctx *config.Context) (*markerTag, error) {
	tpl := tmpl.New(ctx)

	keyVal, err := tpl.Apply(*ctx.Project.Templates.MarkerTagKey)
	if err != nil {
		return nil, err
	}

	valVal, err := tpl.Apply(*ctx.Project.Templates.MarkerTagValue)
	if err != nil {
		return nil, err
	}

	marker := &markerTag{
		key:   keyVal,
		value: valVal,
	}

	return marker, nil
}
