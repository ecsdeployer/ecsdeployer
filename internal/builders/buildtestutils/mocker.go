package buildtestutils

import (
	"fmt"
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
)

func StartMocker(t *testing.T) {

	mocks := []*awsmocker.MockedEndpoint{
		testutil.Mock_ELBv2_DescribeTargetGroups_Generic_Success(),
		testutil.Mock_EC2_DescribeSubnets_Simple(),
		testutil.Mock_EC2_DescribeSecurityGroups_Simple(),
	}

	// add the ones that output to console
	if val := os.Getenv("DEBUG_BUILDERS"); val != "" && testing.Verbose() {
		mocks = append(mocks,
			Mock_ECS_RegisterTaskDefinition_Dump(t),
			Mock_ECS_CreateService_Dump(t),
			Mock_ECS_UpdateService_Dump(t),
		)
	}

	mocks = append(mocks,
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_ECS_RunTask(),
		testutil.Mock_ECS_CreateService_Generic(),
	)

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: mocks,
	})
}

//nolint:unused
func Mock_ECS_RegisterTaskDefinition_Dump(t *testing.T) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "RegisterTaskDefinition",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				prettyJSON, _ := util.JsonifyPretty(rr.JsonPayload)
				t.Log("JSON PAYLOAD:", prettyJSON)

				taskName, _ := jmespath.Search("family", rr.JsonPayload)

				payload, _ := util.Jsonify(map[string]interface{}{
					"taskDefinition": map[string]interface{}{
						"taskDefinitionArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/%s:999", rr.Region, awsmocker.DefaultAccountId, taskName.(string)),
					},
				})

				return payload
			},
		},
	}
}

func Mock_ECS_CreateService_Dump(t *testing.T) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "CreateService",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				prettyJSON, _ := util.JsonifyPretty(rr.JsonPayload)
				t.Log("JSON PAYLOAD:", prettyJSON)

				serviceName, _ := jmespath.Search("serviceName", rr.JsonPayload)
				cluster, _ := jmespath.Search("cluster", rr.JsonPayload)

				payload, _ := util.Jsonify(map[string]interface{}{
					"service": map[string]interface{}{
						"serviceName": serviceName.(string),
						"serviceArn":  fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s/%s", rr.Region, awsmocker.DefaultAccountId, cluster.(string), serviceName.(string)),
					},
				})

				return payload

			},
		},
	}
}

func Mock_ECS_UpdateService_Dump(t *testing.T) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "UpdateService",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				prettyJSON, _ := util.JsonifyPretty(rr.JsonPayload)
				t.Log("JSON PAYLOAD:", prettyJSON)

				serviceName, _ := jmespath.Search("service", rr.JsonPayload)
				cluster, _ := jmespath.Search("cluster", rr.JsonPayload)

				payload, _ := util.Jsonify(map[string]interface{}{
					"service": map[string]interface{}{
						"serviceName": serviceName.(string),
						"serviceArn":  fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s/%s", rr.Region, awsmocker.DefaultAccountId, cluster.(string), serviceName.(string)),
					},
				})

				return payload

			},
		},
	}
}
