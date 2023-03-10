package helpers

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// This will merge multiple arrays of NameValuePairs and then return a unique array
func NameValuePairMerger(pairs ...[]config.NameValuePair) []config.NameValuePair {
	var tempMap map[string]string = make(map[string]string)

	// remove duplicates
	for _, pairGroup := range pairs {

		if pairGroup == nil {
			continue
		}

		for _, pair := range pairGroup {
			if pair.Name == nil || pair.Value == nil {
				continue
			}

			tempMap[*pair.Name] = *pair.Value
		}
	}

	newPairArray := make([]config.NameValuePair, 0, len(tempMap))

	for k, v := range tempMap {

		newPairArray = append(newPairArray, config.NameValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})

	}
	return newPairArray
}

func NameValuePairTemplater(ctx *config.Context, fields tmpl.Fields, pairs []config.NameValuePair) ([]config.NameValuePair, error) {
	result := make([]config.NameValuePair, 0, len(pairs))

	tpl := tmpl.New(ctx).WithExtraFields(fields)

	for _, pair := range pairs {

		keyVal, err := tpl.Apply(*pair.Name)
		if err != nil {
			return nil, err
		}

		valVal, err := tpl.Apply(*pair.Value)
		if err != nil {
			return nil, err
		}

		keyStr := aws.String(keyVal)
		valStr := aws.String(valVal)

		if util.IsBlank(valStr) || util.IsBlank(keyStr) {
			// if either value is blank, that is an invalid tag and we should remove it
			continue
		}

		result = append(result, config.NameValuePair{
			Name:  keyStr,
			Value: valStr,
		})
	}

	return result, nil
}

// build up a list of tags using the various NVPs
func NameValuePair_Build_Tags[T any](ctx *config.Context, thisTags []config.NameValuePair, tplFields tmpl.Fields, buildFunc func(string, string) T, extraTags ...[]config.NameValuePair) ([]T, map[string]string, error) {

	additionalTags := make([]config.NameValuePair, 0, len(extraTags)*5)

	for _, extraNvpList := range extraTags {
		additionalTags = append(additionalTags, extraNvpList...)
	}

	srcTags := NameValuePairMerger(ctx.Project.Tags, ctx.Project.TaskDefaults.Tags, additionalTags, thisTags)
	srcTags, err := NameValuePairTemplater(ctx, tplFields, srcTags)
	if err != nil {
		return nil, nil, err
	}

	tagMap := make(map[string]string, len(srcTags))
	for _, nvp := range srcTags {
		tagMap[*nvp.Name] = *nvp.Value
	}

	if buildFunc == nil {
		return nil, tagMap, nil
	}

	tagList := make([]T, 0, len(tagMap))

	for k, v := range tagMap {
		tagList = append(tagList, buildFunc(k, v))
	}

	return tagList, tagMap, nil
}
