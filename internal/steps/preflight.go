package steps

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

// Run some preflight checks to make sure we can even deploy this thing
func PreflightStep(project *config.Project) *Step {
	return NewStep(&Step{
		Label:    "Preflight",
		Resource: project,
		Create:   stepPreflightCreate,
	})
}

func stepPreflightCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	// VERSION RESTRICTION
	// this is handled in the command

	// ACCOUNT RESTRICTION
	if !ctx.Project.EcsDeployerOptions.IsAllowedAccountId(ctx.AwsAccountId()) {
		return nil, fmt.Errorf("Account '%s' is not an allowed account. Only '%s' is allowed.", ctx.AwsAccountId(), *ctx.Project.EcsDeployerOptions.AllowedAccountId)
	}

	// TEMPLATES
	step.Logger.Debug("Validating Templates")
	if err := stepPreflight_CheckTemplates(ctx); err != nil {
		return nil, err
	}

	// RESOLVE CONTAINER IMAGES
	for _, image := range util.DeepFindInStruct[config.ImageUri](ctx.Project) {
		step.Logger.Debug("Resolving container image")
		if _, err := helpers.ResolveImageUri(ctx, image); err != nil {
			return nil, err
		}
	}

	// CLUSTER
	if ctx.Project.Cluster != nil {
		step.Logger.Debug("Validating Cluster")
		if _, err := ctx.Project.Cluster.Name(ctx); err != nil {
			return nil, err
		}
		if _, err := ctx.Project.Cluster.Arn(ctx); err != nil {
			return nil, err
		}
	}

	// ROLES
	for _, role := range util.DeepFindInStruct[config.RoleArn](ctx.Project) {
		step.Logger.Debug("Validating role")
		if _, err := role.Name(ctx); err != nil {
			return nil, err
		}
		if _, err := role.Arn(ctx); err != nil {
			return nil, err
		}
	}

	// TARGET GROUPS
	for _, tg := range util.DeepFindInStruct[config.TargetGroupArn](ctx.Project) {
		step.Logger.Debug("Validating Target Group")
		if _, err := tg.Name(ctx); err != nil {
			return nil, err
		}
		if _, err := tg.Arn(ctx); err != nil {
			return nil, err
		}
	}

	// NETWORK
	for _, network := range util.DeepFindInStruct[config.NetworkConfiguration](ctx.Project) {
		step.Logger.Debugf("Validating Network")
		if err := network.Resolve(ctx, nil); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func stepPreflight_CheckTemplates(ctx *config.Context) error {
	tpl := tmpl.New(ctx).WithExtraFields(tmpl.Fields{
		"Name":      "THING",
		"Container": "THING",
		"Arch":      "amd64",
	})

	for _, val := range util.DeepFindInStruct[string](ctx.Project.Templates) {
		_, err := tpl.Apply(*val)
		if err != nil {
			return err
		}
	}

	return nil
}
