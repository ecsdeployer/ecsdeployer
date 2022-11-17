package testutil

import (
	"fmt"

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
				return jsonify(map[string]interface{}{
					"RuleArn": fmt.Sprintf("arn:aws:events:%s:%s:rule/%s", rr.Region, awsmocker.DefaultAccountId, name.(string)),
				})
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
				return jsonify(map[string]interface{}{
					"FailedEntries":    []string{},
					"FailedEntryCount": 0,
				})
			},
		},
	}
}

func Mock_Events_ListTargetsByRule(ruleName, busName string, targetIds []string) *awsmocker.MockedEndpoint {
	jmesMatches := map[string]interface{}{
		"Rule": ruleName,
	}

	if busName != "" {
		jmesMatches["EventBusName"] = busName
	}

	results := make([]interface{}, 0, len(targetIds))

	for _, id := range targetIds {
		entry := map[string]interface{}{
			"Id":  id,
			"Arn": fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/testcluster", awsmocker.DefaultRegion, awsmocker.DefaultAccountId),
			"EcsParameters": map[string]interface{}{
				"TaskCount": 1,
			},
		}
		results = append(results, entry)
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "events",
			Action:  "ListTargetsByRule",
			Matcher: JmesRequestMatcher(jmesMatches),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: jsonify(map[string]interface{}{
				"Targets": results,
			}),
		},
	}
}

func Mock_Events_RemoveTargets(ruleName, busName, targetId string) *awsmocker.MockedEndpoint {
	jmesMatches := map[string]interface{}{
		"Rule":   ruleName,
		"Ids[0]": targetId,
	}

	if busName != "" {
		jmesMatches["EventBusName"] = busName
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "events",
			Action:  "RemoveTargets",
			Matcher: JmesRequestMatcher(jmesMatches),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: jsonify(map[string]interface{}{
				"FailedEntries":    []string{},
				"FailedEntryCount": 0,
			}),
		},
	}
}

func Mock_Events_DeleteRule(ruleName, busName string) *awsmocker.MockedEndpoint {
	jmesMatches := map[string]interface{}{
		"Name": ruleName,
	}

	if busName != "" {
		jmesMatches["EventBusName"] = busName
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "events",
			Action:  "DeleteRule",
			Matcher: JmesRequestMatcher(jmesMatches),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			StatusCode:  200,
			Body:        "",
		},
	}
}
