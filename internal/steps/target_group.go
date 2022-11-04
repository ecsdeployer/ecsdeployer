package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func TargetGroupStep(resource *config.Service) *Step {
	return NewStep(&Step{
		Label:    "TargetGroup",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepTargetGroupCreate,
		// Exists:   stepTargetGroupExist,
	})
}

func stepTargetGroupCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {
	return nil, nil
}
