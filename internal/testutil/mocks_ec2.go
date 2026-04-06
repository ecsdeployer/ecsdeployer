package testutil

import (
	"github.com/webdestroya/awsmocker"
)

func Mock_EC2_DescribeSubnets_Simple() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ec2",
			Action:  "DescribeSubnets",
		},
		Response: &awsmocker.MockedResponse{
			DoNotWrap: true,
			Body: map[string]any{
				"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
				"subnetSet": map[string]any{
					"item": []any{
						map[string]any{
							"subnetId": "subnet-633333333333",
							"vpcId":    "vpc-123456789",
						},
						map[string]any{
							"subnetId": "subnet-644444444444",
							"vpcId":    "vpc-123456789",
						},
					},
				},
			},
		},
	}
}

func Mock_EC2_DescribeSecurityGroups_Simple() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ec2",
			Action:  "DescribeSecurityGroups",
		},
		Response: &awsmocker.MockedResponse{
			DoNotWrap: true,
			Body: map[string]any{
				"requestId": "43e9cb52-0e10-40fe-b457-988c8fbfea26",
				"securityGroupInfo": map[string]any{
					"item": []any{
						map[string]any{
							"groupId":   "sg-633333333333",
							"groupName": "fakesg1",
						},
						map[string]any{
							"groupId":   "sg-644444444444",
							"groupName": "fakesg2",
						},
					},
				},
			},
		},
	}
}
