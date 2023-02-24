package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func tagBuilder(k, v string) ecsTypes.Tag {
	return ecsTypes.Tag{
		Key:   &k,
		Value: &v,
	}
}

func PipelineBuild(ctx *config.Context, resource config.IsTaskStruct) (*ecs.RegisterTaskDefinitionInput, error) {

	common, err := config.ExtractCommonTaskAttrs(resource)
	if err != nil {
		return nil, err
	}

	project := ctx.Project
	taskDefaults := project.TaskDefaults
	templates := project.Templates

	arch := util.Coalesce(common.Architecture, taskDefaults.Architecture, util.Ptr(config.ArchitectureDefault))

	commonTplFields, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	if err != nil {
		return nil, err
	}
	tpl := tmpl.New(ctx).WithExtraFields(commonTplFields)

	familyName, err := tpl.Apply(*templates.TaskFamily)
	if err != nil {
		return nil, err
	}

	// TAGS
	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, common.Tags, commonTplFields, tagBuilder)
	if err != nil {
		return nil, err
	}

	taskDef := &ecs.RegisterTaskDefinitionInput{
		NetworkMode:             ecsTypes.NetworkModeAwsvpc,
		Family:                  aws.String(familyName),
		RequiresCompatibilities: []ecsTypes.Compatibility{ecsTypes.CompatibilityFargate},
		Tags:                    tagList,
		RuntimePlatform: &ecsTypes.RuntimePlatform{
			OperatingSystemFamily: ecsTypes.OSFamilyLinux,
			CpuArchitecture:       arch.ToAws(),
		},
	}

	pi := pipeline.NewPipeItem(ctx, taskDef)

	err = pi.Apply(
		&TaskRolesBuilder{},
		&TaskResourcesBuilder{Resource: resource},
		&PrimaryContainerBuilder{Resource: resource, Tpl: tpl},

		// Should be one of the last:
		&ContainerImageResolver{Tpl: tpl},
	)
	if err != nil {
		return nil, err
	}

	return pi.GetData(), nil
}
