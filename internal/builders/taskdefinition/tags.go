package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func tagBuilder(k, v string) ecsTypes.Tag {
	return ecsTypes.Tag{
		Key:   &k,
		Value: &v,
	}
}

func (b *Builder) applyTags() error {

	commonTplFields, err := helpers.GetDefaultTaskTemplateFields(b.ctx, b.commonTask)
	if err != nil {
		return err
	}

	tagList, _, err := helpers.NameValuePair_Build_Tags(b.ctx, b.commonTask.Tags, commonTplFields, tagBuilder)
	if err != nil {
		return err
	}
	b.taskDef.Tags = tagList

	return nil
}
