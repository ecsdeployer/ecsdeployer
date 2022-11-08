package testutil

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
)

func Mock_Events_PutRule_Generic() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "events",
			Action:  "PutRule",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				name, _ := jmespath.Search("Name", rr.JsonPayload)
				return util.Must(util.Jsonify(map[string]interface{}{
					"RuleArn": fmt.Sprintf("arn:aws:events:%s:%s:rule/%s", rr.Region, awsmocker.DefaultAccountId, name.(string)),
				}))
			},
		},
	}
}

func Mock_Events_PutTargets_Generic() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "events",
			Action:  "PutTargets",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				return util.Must(util.Jsonify(map[string]interface{}{
					"FailedEntries":    []string{},
					"FailedEntryCount": 0,
				}))
			},
		},
	}
}
