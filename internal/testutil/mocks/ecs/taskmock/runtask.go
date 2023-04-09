package taskmock

import (
	"fmt"
	"path"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/jmespath/go-jmespath"
	"github.com/webdestroya/awsmocker"
)

func mockRunTask(family, taskId string) *awsmocker.MockedEndpoint {

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "ecs",
			Action:  "RunTask",
			Matcher: func(rr *awsmocker.ReceivedRequest) bool {
				taskDef := testutil.JmesSearchOrNil(rr.JsonPayload, "taskDefinition")
				if taskDef == nil {
					return false
				}

				taskDefArn := taskDef.(string)

				providedFamily := helpers.GetTaskDefFamilyFromArn(taskDefArn)

				return providedFamily == family
			},
		},
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				cluster, _ := jmespath.Search("cluster", rr.JsonPayload)

				clusterName := path.Base(cluster.(string))

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks": []interface{}{
						map[string]interface{}{
							"taskArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task/%s/%s", rr.Region, awsmocker.DefaultAccountId, clusterName, taskId),
						},
					},
				})
			},
		},
	}
}
