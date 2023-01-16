package config

import (
	"net/url"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

type networkTestStruct struct {
	Network *NetworkConfiguration `yaml:"network"`
}

func TestNetwork_Valid(t *testing.T) {
	_, err := yaml.ParseYAMLFile[networkTestStruct]("testdata/network/combined.yml")
	require.NoError(t, err)
}
func TestNetwork_NetworkConfigurationResolver_Errors(t *testing.T) {
	_, _, ctx := networkFilterMocker(t, "testdata/network/filter_test_sgvpc.yml", []*awsmocker.MockedEndpoint{})

	_, err := networkConfigurationResolver(ctx, nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "must provide network")

	_, err = networkConfigurationResolver(ctx, &NetworkConfiguration{})
	require.Error(t, err)
	require.ErrorContains(t, err, "subnets")
	require.ErrorContains(t, err, "No configuration provided")

	_, err = networkConfigurationResolver(ctx, &NetworkConfiguration{
		Subnets: []NetworkFilter{
			util.Must(newNetworkFilterOrIdFromString("subnet-1111")),
		},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "security groups")
	require.ErrorContains(t, err, "No configuration provided")
}

func TestNetwork_NetworkConfigurationResolver_Full(t *testing.T) {
	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSubnets",
				Params: url.Values{
					"Filter.1.Name":    []string{"tag:cloud87/network"},
					"Filter.1.Value.1": []string{"private"},
					"Filter.2.Name":    []string{"state"},
					"Filter.2.Value.1": []string{"available"},
					"Filter.3.Name":    []string{"tag:cloud87/subnet_class"},
					"Filter.3.Value.1": []string{"host"},
				},
			},
			Response: &awsmocker.MockedResponse{
				DoNotWrap: true,
				Body: map[string]interface{}{
					"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
					"subnetSet": map[string]interface{}{
						"item": []interface{}{
							map[string]interface{}{
								"subnetId": "subnet-633333333333",
								"vpcId":    "vpc-123456789",
							},
							map[string]interface{}{
								"subnetId": "subnet-644444444444",
								"vpcId":    "vpc-123456789",
							},
						},
					},
				},
			},
		},
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSecurityGroups",
				Params: url.Values{
					"Filter.1.Name":    []string{"thing"},
					"Filter.1.Value.1": []string{"stuff"},
					"Filter.2.Name":    []string{"test"},
					"Filter.2.Value.1": []string{"test"},
					"Filter.3.Name":    []string{"vpc-id"},
					"Filter.3.Value.1": []string{"vpc-123456789"},
				},
			},
			Response: &awsmocker.MockedResponse{
				DoNotWrap: true,
				Body: map[string]interface{}{
					"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
					"securityGroupInfo": map[string]interface{}{
						"item": []interface{}{
							map[string]interface{}{
								"groupId":   "sg-633333333333",
								"groupName": "fakesg1",
							},
							map[string]interface{}{
								"groupId":   "sg-644444444444",
								"groupName": "fakesg2",
							},
						},
					},
				},
			},
		},
	})

	result, err := networkConfigurationResolver(ctx, network)
	require.NoError(t, err)

	require.NoError(t, network.resolve(ctx))

	ecsNetwork, err := network.ResolveECS(ctx)
	require.NoError(t, err)

	cweNetwork, err := network.ResolveCWE(ctx)
	require.NoError(t, err)

	tables := []struct {
		subnets        []string
		securityGroups []string
		publicIp       bool
	}{
		{result.subnets, result.securityGroups, result.publicIp},
		{ecsNetwork.AwsvpcConfiguration.Subnets, ecsNetwork.AwsvpcConfiguration.SecurityGroups, ecsNetwork.AwsvpcConfiguration.AssignPublicIp == ecsTypes.AssignPublicIpEnabled},
		{cweNetwork.AwsvpcConfiguration.Subnets, cweNetwork.AwsvpcConfiguration.SecurityGroups, cweNetwork.AwsvpcConfiguration.AssignPublicIp == eventTypes.AssignPublicIpEnabled},
	}

	for _, table := range tables {
		require.Contains(t, table.subnets, "subnet-00000000000")
		require.Contains(t, table.subnets, "subnet-11111111111")
		require.Contains(t, table.subnets, "subnet-633333333333")
		require.Contains(t, table.subnets, "subnet-644444444444")
		require.Len(t, table.subnets, 4)

		require.Contains(t, table.securityGroups, "sg-1234567890")
		require.Contains(t, table.securityGroups, "sg-9876543210")
		require.Contains(t, table.securityGroups, "sg-633333333333")
		require.Contains(t, table.securityGroups, "sg-644444444444")
		require.Len(t, table.securityGroups, 4)

		require.False(t, table.publicIp)

	}

}

func TestNetworkCalculateSecurityGroups_NoVPC(t *testing.T) {
	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test_nofilter.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSecurityGroups",
			},
			Response: &awsmocker.MockedResponse{
				StatusCode: 400,
				Body:       "EXPLODE",
			},
		},
	})

	result, err := calculateSecurityGroups(ctx, *network, nil)
	require.NoError(t, err)

	require.Contains(t, result, "sg-1234567890")
	require.Contains(t, result, "sg-9876543210")
	require.Len(t, result, 2)
}

func TestNetworkCalculateSecurityGroups_WithVPC(t *testing.T) {

	vpcId := "vpc-11111111"

	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSecurityGroups",
				Params: url.Values{
					"Filter.1.Name":    []string{"thing"},
					"Filter.1.Value.1": []string{"stuff"},
					"Filter.2.Name":    []string{"test"},
					"Filter.2.Value.1": []string{"test"},
					"Filter.3.Name":    []string{"vpc-id"},
					"Filter.3.Value.1": []string{vpcId},
				},
			},
			Response: &awsmocker.MockedResponse{
				DoNotWrap: true,
				Body: map[string]interface{}{
					"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
					"securityGroupInfo": map[string]interface{}{
						"item": []interface{}{
							map[string]interface{}{
								"groupId":   "sg-33333333333",
								"groupName": "fakesg1",
								"vpcId":     vpcId,
							},
							map[string]interface{}{
								"groupId":   "sg-44444444444",
								"groupName": "fakesg2",
								"vpcId":     vpcId,
							},
						},
					},
				},
			},
		},
	})

	result, err := calculateSecurityGroups(ctx, *network, &vpcId)
	require.NoError(t, err)

	require.Contains(t, result, "sg-1234567890")
	require.Contains(t, result, "sg-9876543210")
	require.Contains(t, result, "sg-33333333333")
	require.Contains(t, result, "sg-44444444444")
	require.Len(t, result, 4)
}

// if the filter list already has a VPC filter, don't add our own
func TestNetworkCalculateSecurityGroups_WithExistingVPC(t *testing.T) {

	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test_sgvpc.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSecurityGroups",
				Params: url.Values{
					"Filter.1.Name":    []string{"vpc-id"},
					"Filter.1.Value.1": []string{"vpc-999999999"},
					"Filter.2.Name":    []string{"test"},
					"Filter.2.Value.1": []string{"test"},
				},
			},
			Response: &awsmocker.MockedResponse{
				DoNotWrap: true,
				Body: map[string]interface{}{
					"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
					"securityGroupInfo": map[string]interface{}{
						"item": []interface{}{
							map[string]interface{}{
								"groupId":   "sg-533333333333",
								"groupName": "fakesg1",
							},
							map[string]interface{}{
								"groupId":   "sg-544444444444",
								"groupName": "fakesg2",
							},
						},
					},
				},
			},
		},
	})

	result, err := calculateSecurityGroups(ctx, *network, util.Ptr("vpc-11111111"))
	require.NoError(t, err)

	require.Contains(t, result, "sg-1234567890")
	require.Contains(t, result, "sg-9876543210")
	require.Contains(t, result, "sg-533333333333")
	require.Contains(t, result, "sg-544444444444")
	require.Len(t, result, 4)
}

func TestNetworkCalculateSubnets_NoVPC(t *testing.T) {
	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test_nofilter.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSubnets",
			},
			Response: &awsmocker.MockedResponse{
				StatusCode: 400,
				Body:       "EXPLODE",
			},
		},
	})

	result, vpcid, err := calculateSubnets(ctx, *network)
	require.NoError(t, err)

	require.Nil(t, vpcid, "expected VPCid to be nil")

	require.Contains(t, result, "subnet-00000000000")
	require.Contains(t, result, "subnet-11111111111")
	require.Len(t, result, 2)
}

func TestNetworkCalculateSubnets_WithVPC(t *testing.T) {
	network, _, ctx := networkFilterMocker(t, "testdata/network/filter_test.yml", []*awsmocker.MockedEndpoint{
		{
			Request: &awsmocker.MockedRequest{
				Service: "ec2",
				Action:  "DescribeSubnets",
				Params: url.Values{
					"Filter.1.Name":    []string{"tag:cloud87/network"},
					"Filter.1.Value.1": []string{"private"},
					"Filter.2.Name":    []string{"state"},
					"Filter.2.Value.1": []string{"available"},
					"Filter.3.Name":    []string{"tag:cloud87/subnet_class"},
					"Filter.3.Value.1": []string{"host"},
				},
			},
			Response: &awsmocker.MockedResponse{
				DoNotWrap: true,
				Body: map[string]interface{}{
					"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
					"subnetSet": map[string]interface{}{
						"item": []interface{}{
							map[string]interface{}{
								"subnetId": "subnet-33333333333",
								"vpcId":    "vpc-123456789",
							},
							map[string]interface{}{
								"subnetId": "subnet-44444444444",
								"vpcId":    "vpc-123456789",
							},
						},
					},
				},
			},
		},
	})

	result, vpcid, err := calculateSubnets(ctx, *network)
	require.NoError(t, err)

	require.NotNil(t, vpcid, "expected VPCid to not be nil")
	require.EqualValues(t, "vpc-123456789", *vpcid)

	require.Contains(t, result, "subnet-00000000000")
	require.Contains(t, result, "subnet-11111111111")
	require.Contains(t, result, "subnet-33333333333")
	require.Contains(t, result, "subnet-44444444444")
	require.Len(t, result, 4)

}

func networkFilterMocker(t *testing.T, filePath string, mocks []*awsmocker.MockedEndpoint) (*NetworkConfiguration, *Project, *Context) {
	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: append([]*awsmocker.MockedEndpoint{}, mocks...),
	})

	network, err := yaml.ParseYAMLFile[NetworkConfiguration](filePath)
	require.NoError(t, err)

	project, err := yaml.ParseYAMLFile[Project]("testdata/simple.yml")
	require.NoError(t, err)

	ctx := New(project)

	return network, project, ctx
}
