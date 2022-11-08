package testutil

import (
	"fmt"
	"path"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
)

func Mock_ECS_RunTask() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "RunTask",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				cluster, _ := jmespath.Search("cluster", rr.JsonPayload)

				clusterName := path.Base(cluster.(string))

				return util.Must(util.Jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks": []interface{}{
						map[string]interface{}{
							"taskArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task/%s/deadbeefdeadbeefdeadbeefdeadbeef", rr.Region, awsmocker.DefaultAccountId, clusterName),
						},
					},
				}))
			},
		},
	}
}

func Mock_ECS_DescribeTasks_Pending(status string, maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeTasks",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				clusterName := path.Base(jmespathSearch(rr.JsonPayload, "cluster").(string))
				taskArn := jmespathSearch(rr.JsonPayload, "tasks[0]").(string)

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks": []interface{}{
						map[string]interface{}{
							"lastStatus": status,
							"clusterArn": fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", rr.Region, awsmocker.DefaultAccountId, clusterName),
							"taskArn":    taskArn,
						},
					},
				})
			},
		},
	}
}

func Mock_ECS_DescribeTasks_Stopped(status string, maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeTasks",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				clusterName := path.Base(jmespathSearch(rr.JsonPayload, "cluster").(string))
				taskArn := jmespathSearch(rr.JsonPayload, "tasks[0]").(string)

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks": []interface{}{
						map[string]interface{}{
							"lastStatus": status,
							"clusterArn": fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", rr.Region, awsmocker.DefaultAccountId, clusterName),
							"taskArn":    taskArn,
						},
					},
				})
			},
		},
	}
}

func Mock_ECS_RegisterTaskDefinition_Generic() *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "RegisterTaskDefinition",
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				taskName, _ := jmespath.Search("family", rr.JsonPayload)

				return util.Must(util.Jsonify(map[string]interface{}{
					"taskDefinition": map[string]interface{}{
						"taskDefinitionArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/%s:999", rr.Region, awsmocker.DefaultAccountId, taskName.(string)),
					},
				}))
			},
		},
	}
}
