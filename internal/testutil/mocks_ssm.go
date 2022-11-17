package testutil

import (
	"fmt"
	"time"

	"github.com/webdestroya/awsmocker"
	"golang.org/x/exp/slices"
)

func Mock_SSM_GetParametersByPath(prefixWithTrailingSlash string, paramNames []string) *awsmocker.MockedEndpoint {
	jmesMatches := map[string]interface{}{
		"Path":           prefixWithTrailingSlash,
		"Recursive":      true,
		"WithDecryption": false,
	}

	results := make([]interface{}, 0, len(paramNames))

	slices.Sort(paramNames)

	for _, paramName := range paramNames {

		name := prefixWithTrailingSlash + paramName

		entry := map[string]interface{}{
			"Name":             name,
			"Type":             "SecureString",
			"Version":          1,
			"DataType":         "text",
			"ARN":              fmt.Sprintf("arn:aws:ssm:%s:%s:parameter%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, name),
			"LastModifiedDate": time.Now().UTC().Unix(),
			"Value":            "ZmFrZQ==",
		}

		results = append(results, entry)
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ssm",
			Action:  "GetParametersByPath",
			Matcher: JmesRequestMatcher(jmesMatches),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: jsonify(map[string]interface{}{
				"Parameters": results,
			}),
		},
	}
}
