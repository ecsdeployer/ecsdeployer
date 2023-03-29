package testutil

import (
	"fmt"
	"time"

	"github.com/webdestroya/awsmocker"
)

func Mock_Logs_CreateLogGroup_AllowAny() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "CreateLogGroup",
		},
		Response: MockResponse_EmptySuccess(),
	}
}

func Mock_Logs_CreateLogGroup(logGroupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "CreateLogGroup",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupName": logGroupName,
			}),
		},
		Response: MockResponse_EmptySuccess(),
	}
}

func Mock_Logs_CreateLogGroup_Deny(logGroupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "CreateLogGroup",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupName": logGroupName,
			}),
		},
		Response: awsmocker.MockResponse_Error(400, "LimitExceededException", "You have reached the maximum number of resources that can be created."),
	}
}

func Mock_Logs_CreateLogGroup_AlreadyExists(logGroupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "CreateLogGroup",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupName": logGroupName,
			}),
		},
		Response: awsmocker.MockResponse_Error(400, "ResourceAlreadyExistsException", "The specified resource already exists."),
	}
}

func Mock_Logs_PutRetentionPolicy_AllowAny() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "PutRetentionPolicy",
		},
		Response: MockResponse_EmptySuccess(),
	}
}
func Mock_Logs_PutRetentionPolicy(logGroupName string, days int32) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "PutRetentionPolicy",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupName":    logGroupName,
				"retentionInDays": days,
			}),
		},
		Response: MockResponse_EmptySuccess(),
	}
}

func Mock_Logs_DeleteRetentionPolicy(logGroupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "DeleteRetentionPolicy",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupName": logGroupName,
			}),
		},
		Response: MockResponse_EmptySuccess(),
	}
}

func Mock_Logs_DescribeLogGroups(logGroupRetentions map[string]int32) *awsmocker.MockedEndpoint {

	if logGroupRetentions == nil {
		logGroupRetentions = make(map[string]int32)
	}

	results := make([]interface{}, 0, len(logGroupRetentions))
	for k, v := range logGroupRetentions {

		entry := map[string]interface{}{
			"arn":               fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, k),
			"logGroupName":      k,
			"storedBytes":       0,
			"creationTime":      time.Now().AddDate(0, -1, 0).UTC().Unix(),
			"metricFilterCount": 0,
		}
		if v > 0 {
			entry["retentionInDays"] = v
		}

		results = append(results, entry)
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "DescribeLogGroups",
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: jsonify(map[string]interface{}{
				"logGroups": results,
			}),
		},
	}
}

func Mock_Logs_DescribeLogGroups_Single(logGroupName string, retention int32) *awsmocker.MockedEndpoint {

	entry := map[string]interface{}{
		"arn":               fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, logGroupName),
		"logGroupName":      logGroupName,
		"storedBytes":       0,
		"creationTime":      time.Now().AddDate(0, -1, 0).UTC().Unix(),
		"metricFilterCount": 0,
	}
	if retention > 0 {
		entry["retentionInDays"] = retention
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "DescribeLogGroups",
			Matcher: JmesRequestMatcher(map[string]interface{}{
				"logGroupNamePrefix": logGroupName,
			}),
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeJSON,
			Body: jsonify(map[string]interface{}{
				"logGroups": []interface{}{entry},
			}),
		},
	}
}
