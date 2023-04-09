package testutil

import (
	"fmt"
	"sync"
	"time"

	"github.com/webdestroya/awsmocker"
	"golang.org/x/exp/slices"
)

func Mock_SSM_GetParametersByPath(prefixWithTrailingSlash string, paramNames []string) *awsmocker.MockedEndpoint {
	return Mock_SSM_GetParametersByPath_Advanced(func(m *Mock_ECS_GetParametersByPathOpts) {
		m.Path = prefixWithTrailingSlash
		m.Names = paramNames
	})
}

type Mock_ECS_GetParametersByPathOpts struct {
	MaxCount  int
	Path      string
	Names     []string
	NumParams int
	NextToken bool
}

var (
	ssmGenCount int
	ssmGenMu    sync.Mutex
)

func Mock_SSM_GetParametersByPath_Advanced(optFuncs ...func(*Mock_ECS_GetParametersByPathOpts)) *awsmocker.MockedEndpoint {

	options := Mock_ECS_GetParametersByPathOpts{}
	for _, optFunc := range optFuncs {
		optFunc(&options)
	}

	jmesMatches := map[string]interface{}{
		"Path":           options.Path,
		"Recursive":      true,
		"WithDecryption": false,
	}

	// response := map[string]interface{}{
	// 	"Parameters": results,
	// }

	// if options.NextToken {
	// 	response["NextToken"] = RandomHex(32)
	// }

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ssm",
			Action:        "GetParametersByPath",
			MaxMatchCount: options.MaxCount,
			Matcher:       JmesRequestMatcher(jmesMatches),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: func(rr *awsmocker.ReceivedRequest) string {

				params := slices.Clone(options.Names)

				if options.NumParams > 0 {
					if len(params) == 0 {
						params = make([]string, 0, options.NumParams)
					}
					ssmGenMu.Lock()
					defer ssmGenMu.Unlock()
					for i := 0; i < options.NumParams; i++ {
						ssmGenCount++
						params = append(params, fmt.Sprintf("SSM_VAR_F%02d", ssmGenCount))
					}
				}

				results := make([]interface{}, 0, len(params))

				slices.Sort(params)

				for _, paramName := range params {

					name := options.Path + paramName

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

				response := map[string]interface{}{
					"Parameters": results,
				}

				if options.NextToken {
					response["NextToken"] = RandomHex(32)
				}

				return jsonify(response)
			},
			// Body:        func(jsonify(response),
		},
	}
}
