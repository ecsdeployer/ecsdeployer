package testutil

import (
	"github.com/webdestroya/awsmocker"
)

func Mock_Logs_CreateLogGroup() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "CreateLogGroup",
		},
		Response: MockResponse_EmptySuccess(),
	}
}

func Mock_Logs_PutRetentionPolicy() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "logs",
			Action:  "PutRetentionPolicy",
		},
		Response: MockResponse_EmptySuccess(),
	}
}
