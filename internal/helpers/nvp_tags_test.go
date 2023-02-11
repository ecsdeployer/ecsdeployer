package helpers_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestNameValuePair_Build_Tags(t *testing.T) {

	type easyTag struct {
		Name  string
		Value string
	}

	project := &config.Project{
		ProjectName: "dummy",
		Tags: []config.NameValuePair{
			{
				Name:  aws.String("test/{{.ProjectName}}/thing"),
				Value: aws.String("test/{{.Project}}/stuff"),
			},
			{
				Name:  aws.String("projectlevel"),
				Value: aws.String("originalValue"),
			},
			{
				Name:  aws.String("projectlevel2"),
				Value: aws.String("originalValue"),
			},
			{
				Name:  aws.String("projName"),
				Value: aws.String("{{.ProjectName}}"),
			},
		},
		TaskDefaults: &config.FargateDefaults{
			CommonTaskAttrs: config.CommonTaskAttrs{
				Tags: []config.NameValuePair{
					{
						Name:  aws.String("taskdef/{{.ProjectName}}/thing"),
						Value: aws.String("something"),
					},
					{
						Name:  aws.String("projectlevel"),
						Value: aws.String("overridden"),
					},
					{
						Name:  aws.String("projectlevel2"),
						Value: aws.String("taskdef"),
					},
				},
			},
		},
	}
	project.ApplyDefaults()

	ctx := config.New(project)

	fields := tmpl.Fields{}

	thisTags := []config.NameValuePair{
		{
			Name:  aws.String("projectlevel2"),
			Value: aws.String("thisTags"),
		},
	}

	tagList, tagMap, err := helpers.NameValuePair_Build_Tags(ctx, thisTags, fields, func(s1, s2 string) easyTag {
		return easyTag{
			Name:  s1,
			Value: s2,
		}
	})

	require.NoError(t, err)
	require.NotNil(t, tagList)
	require.NotNil(t, tagMap)

	require.Equal(t, tagMap["projectlevel"], "overridden")
	require.Equal(t, tagMap["projName"], "dummy")
	require.Equal(t, tagMap["projectlevel2"], "thisTags")

	for _, pair := range tagList {
		require.Contains(t, tagMap, pair.Name)
		require.Equal(t, tagMap[pair.Name], pair.Value)
	}

}
