package testutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
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

func Mock_Scheduler_GetSchedule_Missing(groupName, scheduleName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("/schedules/%s", scheduleName),
			Params: url.Values{
				"groupName": []string{groupName},
			},
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

func Mock_Scheduler_CreateSchedule(groupName, scheduleName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodPost,
			Path:    fmt.Sprintf("/schedules/%s", scheduleName),
			JMESPathMatches: map[string]interface{}{
				"GroupName": groupName,
			},
		},
		Response: &awsmocker.MockedResponse{
			Body: jsonify(map[string]interface{}{
				"ScheduleArn": fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule/%s/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, groupName, scheduleName),
			}),
		},
	}
}

func Mock_Scheduler_UpdateSchedule(groupName, scheduleName string) *awsmocker.MockedEndpoint {
	endpoint := Mock_Scheduler_CreateSchedule(groupName, scheduleName)
	endpoint.Request.Method = http.MethodPut
	return endpoint
}

func Mock_Scheduler_DeleteSchedule(groupName, scheduleName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodDelete,
			Path:    fmt.Sprintf("/schedules/%s", scheduleName),
			Params: url.Values{
				"groupName": []string{groupName},
			},
		},
		Response: &awsmocker.MockedResponse{
			Body: "",
		},
	}
}

type MockListScheduleEntry struct {
	Name      string
	TargetArn string
	State     schedulerTypes.ScheduleState
}

func (m MockListScheduleEntry) toMap(groupName string) map[string]any {
	result := map[string]any{
		"Arn":                  fmt.Sprintf("arn:aws:scheduler:%s:%s:schedule/%s/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, groupName, m.Name),
		"CreationDate":         time.Now().UnixMilli(),
		"LastModificationDate": time.Now().UnixMilli(),
		"Name":                 m.Name,
		"GroupName":            groupName,
		"State":                m.State,

		"Target": map[string]any{
			"Arn": m.TargetArn,
		},
	}

	return result
}

func Mock_Scheduler_ListSchedules(groupName string, schedules []MockListScheduleEntry) *awsmocker.MockedEndpoint {

	listOfScheds := make([]map[string]any, 0, len(schedules))
	for _, entry := range schedules {
		listOfScheds = append(listOfScheds, entry.toMap(groupName))
	}

	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "scheduler",
			Method:  http.MethodPost,
			Path:    "/schedules",
			Params: url.Values{
				"ScheduleGroup": []string{groupName},
			},
		},
		Response: &awsmocker.MockedResponse{
			Body: jsonify(map[string]interface{}{
				"Schedules": listOfScheds,
			}),
		},
	}
}
