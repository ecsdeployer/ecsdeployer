package service

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

	tagList, _, err := helpers.NameValuePair_Build_Tags(b.ctx, b.entity.Tags, b.commonTplFields, tagBuilder)
	if err != nil {
		return err
	}
	b.serviceDef.Tags = tagList

	return nil
}
