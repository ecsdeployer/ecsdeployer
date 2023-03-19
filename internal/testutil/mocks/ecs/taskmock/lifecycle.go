package taskmock

import (
	"fmt"
	"path"
	"path/filepath"
	"sync"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/webdestroya/awsmocker"
)

type lifecycle struct {
	taskId    string
	options   *Options
	callCount int
	mu        sync.Mutex

	muStopped  sync.Once
	_stoppedAt int64
}

func (lc *lifecycle) stoppedAt() int64 {
	lc.muStopped.Do(func() {
		lc._stoppedAt = time.Now().UTC().Unix()
	})
	return lc._stoppedAt
}

func (lc *lifecycle) describeTaskMock() *awsmocker.MockedEndpoint {

	runningCallStart := lc.options.PendingCount
	stoppedCallStart := runningCallStart + lc.options.RunningCount

	return &awsmocker.MockedEndpoint{
		Request: describeTaskRequest(lc.taskId),
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {
				lc.mu.Lock()
				lc.callCount++
				defer lc.mu.Unlock()

				clusterName := path.Base(testutil.JmesPathSearch(rr.JsonPayload, "cluster").(string))
				taskArn := testutil.JmesPathSearch(rr.JsonPayload, "tasks[0]").(string)
				clusterArn := fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", rr.Region, awsmocker.DefaultAccountId, clusterName)

				taskResult := map[string]any{
					"taskArn":       taskArn,
					"clusterArn":    clusterArn,
					"desiredStatus": "RUNNING",
					"lastStatus":    "PENDING",
				}

				if lc.callCount > stoppedCallStart {
					taskResult["lastStatus"] = "STOPPED"
					taskResult["desiredStatus"] = "STOPPED"
					taskResult["stoppedAt"] = lc.stoppedAt()
					taskResult["stopCode"] = lc.options.StopReason
					taskResult["stoppedReason"] = "aws stoppedReason here"
					taskResult["containers"] = []interface{}{
						map[string]interface{}{
							"name":     "primary",
							"exitCode": lc.options.ExitCode,
						},
					}
				} else if lc.callCount > runningCallStart {
					taskResult["lastStatus"] = "RUNNING"
				} else {
					// pending
				}

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"tasks":    []interface{}{taskResult},
				})
			},
		},
	}
}

func describeTaskRequest(taskId string) *awsmocker.MockedRequest {
	return &awsmocker.MockedRequest{
		Service: "ecs",
		Action:  "DescribeTasks",
		Matcher: func(rr *awsmocker.ReceivedRequest) bool {
			taskArnRaw := testutil.JmesSearchOrNil(rr.JsonPayload, "tasks[0]")
			if taskArnRaw == nil {
				return false
			}

			taskArn := taskArnRaw.(string)
			return taskId == filepath.Base(taskArn)
		},
	}
}
