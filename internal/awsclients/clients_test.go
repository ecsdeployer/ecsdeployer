package awsclients_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"

	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func TestClientGetters(t *testing.T) {
	testutil.StartMocker(t, nil)

	require.IsType(t, &sts.Client{}, awsclients.STSClient())
	require.IsType(t, &ecs.Client{}, awsclients.ECSClient())
	require.IsType(t, &ssm.Client{}, awsclients.SSMClient())
	require.IsType(t, &ec2.Client{}, awsclients.EC2Client())
	require.IsType(t, &elbv2.Client{}, awsclients.ELBv2Client())
	require.IsType(t, &events.Client{}, awsclients.EventsClient())
	require.IsType(t, &logs.Client{}, awsclients.LogsClient())
	require.IsType(t, &tagging.Client{}, awsclients.TaggingClient())
}
