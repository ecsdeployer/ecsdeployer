package testutil

import (
	"fmt"
	"path"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

func Mock_ECS_RunTask_FailToLaunch(maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "RunTask",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				return util.Must(util.Jsonify(map[string]interface{}{
					"failures": []interface{}{
						map[string]interface{}{
							"detail": "some failure detail",
							"reason": "you goofed",
						},
					},
					"tasks": []interface{}{},
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
				clusterName := path.Base(JmesPathSearch(rr.JsonPayload, "cluster").(string))
				taskArn := JmesPathSearch(rr.JsonPayload, "tasks[0]").(string)

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

func Mock_ECS_DescribeTasks_Stopped(stopCode ecsTypes.TaskStopCode, exitCode int, maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeTasks",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				clusterName := path.Base(JmesPathSearch(rr.JsonPayload, "cluster").(string))
				taskArn := JmesPathSearch(rr.JsonPayload, "tasks[0]").(string)

				taskInfo := map[string]interface{}{
					"lastStatus":    "STOPPED",
					"desiredStatus": "STOPPED",
					"clusterArn":    fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", rr.Region, awsmocker.DefaultAccountId, clusterName),
					"taskArn":       taskArn,
					"stoppedAt":     time.Now().UTC().Unix(),
					"stopCode":      stopCode,
					"stoppedReason": "something something",
					"containers": []interface{}{
						map[string]interface{}{
							"name":     "primary",
							"exitCode": exitCode,
						},
					},
				}

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks":    []interface{}{taskInfo},
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

func Mock_ECS_DescribeServices_Missing(maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeServices",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				serviceName := JmesPathSearch(rr.JsonPayload, "services[0]").(string)

				return jsonify(map[string]interface{}{
					"failures": []interface{}{
						map[string]interface{}{
							"arn":    fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, serviceName),
							"reason": "MISSING",
						},
					},
					"services": []interface{}{},
				})
			},
		},
	}
}

func Mock_ECS_DescribeServices_Single(svc ecsTypes.Service, maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeServices",
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				clusterArn := JmesPathSearch(rr.JsonPayload, "cluster").(string)
				serviceName := JmesPathSearch(rr.JsonPayload, "services[0]").(string)

				svc.ClusterArn = aws.String(clusterArn)
				svc.ServiceName = aws.String(serviceName)

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"services": []interface{}{MockResponse_ECS_Service(svc)},
				})
			},
		},
	}
}

func Mock_ECS_DescribeServices_jmespath(jmesMatches map[string]any, svc ecsTypes.Service, maxCount int) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service:       "ecs",
			Action:        "DescribeServices",
			Matcher:       JmesRequestMatcher(jmesMatches),
			MaxMatchCount: maxCount,
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				clusterArn := JmesPathSearch(rr.JsonPayload, "cluster").(string)
				serviceName := JmesPathSearch(rr.JsonPayload, "services[0]").(string)

				svc.ClusterArn = aws.String(clusterArn)
				svc.ServiceName = aws.String(serviceName)

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"services": []interface{}{MockResponse_ECS_Service(svc)},
				})
			},
		},
	}
}

// func Mock_ECS_CreateService() *awsmocker.MockedEndpoint {
// 	return &awsmocker.MockedEndpoint{
// 		Request: &awsmocker.MockedRequest{
// 			Service: "ecs",
// 			Action:  "CreateService",
// 		},
// 		Response: &awsmocker.MockedResponse{
// 			Body: func(rr *awsmocker.ReceivedRequest) string {

// 				return jsonify(map[string]interface{}{
// 					"taskDefinition": map[string]interface{}{
// 						"taskDefinitionArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/%s:999", rr.Region, awsmocker.DefaultAccountId, "BLAH"),
// 					},
// 				})
// 			},
// 		},
// 	}
// }

func Mock_ECS_CreateService_jmespath(jmesMatchers map[string]any, service ecsTypes.Service) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "CreateService",
			Matcher: JmesRequestMatcher(jmesMatchers),
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				return jsonify(map[string]interface{}{
					"service": MockResponse_ECS_Service(service),
				})
			},
		},
	}
}

func Mock_ECS_UpdateService_jmespath(jmesMatchers map[string]any, service ecsTypes.Service) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "UpdateService",
			Matcher: JmesRequestMatcher(jmesMatchers),
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				return jsonify(map[string]interface{}{
					"service": MockResponse_ECS_Service(service),
				})
			},
		},
	}
}

func MockResponse_ECS_Service(service ecsTypes.Service) map[string]interface{} {
	val := make(map[string]interface{}, 20)

	timeNow := time.Now().UTC()
	// val["createdAt"] = timeNow.Format(time.RFC3339) // would be nice if AWS was consistent, but alas... no.
	val["createdAt"] = timeNow.Unix()
	val["platformVersion"] = "LATEST"
	val["platformFamily"] = "Linux"

	if service.ClusterArn != nil {
		val["clusterArn"] = *service.ClusterArn
	}

	val["desiredCount"] = service.DesiredCount // used in waiter
	val["runningCount"] = service.RunningCount // used in waiter
	val["pendingCount"] = service.PendingCount
	val["launchType"] = service.LaunchType

	if service.TaskDefinition != nil {
		val["taskDefinition"] = *service.TaskDefinition
	}

	if service.ServiceArn != nil {
		val["serviceArn"] = *service.ServiceArn
	}
	if service.ServiceName != nil {
		val["serviceName"] = *service.ServiceName
	}
	if service.Status != nil {
		// used in stable waiter
		val["status"] = *service.Status
	}

	// used in waiter
	deps := make([]interface{}, 0, len(service.Deployments))
	for i, dep := range service.Deployments {
		depObj := map[string]interface{}{
			"id":           fmt.Sprintf("ecs-svc/%d", i),
			"updatedAt":    timeNow.Unix(),
			"createdAt":    timeNow.Unix(),
			"desiredCount": dep.DesiredCount,
			"pendingCount": dep.PendingCount,
			"runningCount": dep.RunningCount,
		}

		if dep.Status != nil {
			depObj["status"] = *dep.Status
		}

		if dep.TaskDefinition != nil {
			depObj["taskDefinition"] = *dep.TaskDefinition
		}

		deps = append(deps, depObj)
	}
	val["deployments"] = deps

	if service.ClusterArn != nil {
		val["clusterArn"] = *service.ClusterArn
	}

	return val
}
