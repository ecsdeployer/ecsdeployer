package awsclients_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestAwsConfig(t *testing.T) {
	testutil.StartMocker(t, nil)

	require.IsType(t, aws.Config{}, awsclients.AwsConfig())
}
