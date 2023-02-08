package testutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/webdestroya/awsmocker"
)

func Mock_Scheduler_GetScheduleGroup_Missing(groupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("/schedule-groups/%s", groupName),
		},
		Response: &awsmocker.MockedResponse{
			Handler: func(rr *awsmocker.ReceivedRequest) *http.Response {

				body := jsonify(map[string]interface{}{
					"Message": fmt.Sprintf("Schedule %s does not exist.", groupName),
				})

				status := http.StatusNotFound

				resp := &http.Response{
					StatusCode: status,
					Status:     http.StatusText(status),
					Header:     make(http.Header),
				}

				resp.Header.Set("x-amzn-ErrorType", "ResourceNotFoundException:http://internal.amazon.com/coral/com.amazonaws.chronos/")
				resp.Header.Set("Content-Type", "application/json")

				buf := bytes.NewBufferString(body)

				resp.ContentLength = int64(buf.Len())
				resp.Body = io.NopCloser(buf)

				return resp
			},
		},
	}
}

func Mock_Scheduler_GetScheduleGroup(groupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("/schedule-groups/%s", groupName),
		},
		Response: &awsmocker.MockedResponse{
			Body: jsonify(map[string]interface{}{
				"Arn":                  fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule-group/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, groupName),
				"CreationDate":         time.Now().UnixMilli(),
				"LastModificationDate": time.Now().UnixMilli(),
				"Name":                 groupName,
				"State":                "ACTIVE",
			}),
		},
	}
}

func Mock_Scheduler_CreateScheduleGroup(groupName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodPost,
			Path:    fmt.Sprintf("/schedule-groups/%s", groupName),
		},
		Response: &awsmocker.MockedResponse{
			Body: jsonify(map[string]interface{}{
				"ScheduleGroupArn": fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule-group/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, groupName),
			}),
		},
	}
}
