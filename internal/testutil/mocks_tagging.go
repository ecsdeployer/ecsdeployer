package testutil

import (
	"fmt"

	"github.com/webdestroya/awsmocker"
	"golang.org/x/exp/maps"
)

// resourceArns can be either a string array, or a map or key=>map[tagkey]=tagvalue
func Mock_Tagging_GetResources(typeFilter string, tags map[string]string, resourceArns any) *awsmocker.MockedEndpoint {

	// TagFilters | [?Key=='ecsdeployer/project'].Values[0] | [0]
	jmesMatchers := map[string]interface{}{
		"ResourceTypeFilters[0]": typeFilter,
	}
	for k, v := range tags {
		tagPath := fmt.Sprintf("TagFilters | [?Key=='%s'].Values[0] | [0]", k)
		jmesMatchers[tagPath] = v
	}

	type resMapping struct {
		arn  string
		tags map[string]string
	}

	mappings := make([]resMapping, 0)

	switch rarns := resourceArns.(type) {
	case map[string]map[string]string:

		for rArn, extraTags := range rarns {
			extraTags := extraTags
			maps.Copy(extraTags, tags)

			mappings = append(mappings, resMapping{
				arn:  rArn,
				tags: extraTags,
			})
		}

	case []string:

		for _, rArn := range rarns {
			mappings = append(mappings, resMapping{
				arn:  rArn,
				tags: tags,
			})
		}

	default:
		panic(fmt.Errorf("resourceArns must be either a []string or map[string]map[string]string. You provided %T", resourceArns))
	}

	results := make([]interface{}, 0, len(mappings))

	for _, mapping := range mappings {
		tagList := make([]interface{}, 0, len(mapping.tags))
		for k, v := range mapping.tags {
			tagList = append(tagList, map[string]string{
				"Key":   k,
				"Value": v,
			})
		}
		results = append(results, map[string]interface{}{
			"ResourceARN": mapping.arn,
			"Tags":        tagList,
		})
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "tagging",
			Action:  "GetResources",
			Matcher: JmesRequestMatcher(jmesMatchers),
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				return jsonify(map[string]interface{}{
					// "PaginationToken": "",
					"ResourceTagMappingList": results,
				})
			},
		},
	}
}
