package testutil

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"github.com/webdestroya/awsmocker"
)

func StartMocker(t *testing.T, opts *awsmocker.MockerOptions) *awsmocker.MockerInfo {

	if opts == nil {
		opts = &awsmocker.MockerOptions{}
	}

	// force the mocker to not mess with system proxy
	opts.ReturnAwsConfig = true

	info := awsmocker.Start(t, opts)

	awsclients.SetupWithConfig(*info.AwsConfig)

	return info
}

func MockResponse_EmptySuccess() *awsmocker.MockedResponse {
	return &awsmocker.MockedResponse{
		StatusCode: 200,
		Body:       "OK",
	}
}

// This is just a basic mock server to get the account ID and region
func MockSimpleStsProxy(t *testing.T) {
	StartMocker(t, nil)
}
