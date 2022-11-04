package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func FirelensStep(resource *config.Project) *Step {
	return NewStep(&Step{
		Label:    "Firelens",
		Create:   stepFirelensCreate,
		PreApply: stepFirelensPreApply,
	})
}

func stepFirelensPreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {

	return nil
}

func stepFirelensCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	// Create firelens group

	return nil, nil
}
