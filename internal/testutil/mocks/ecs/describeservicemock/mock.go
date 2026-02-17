package describeservicemock

import (
	"fmt"
	"path/filepath"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/webdestroya/awsmocker"
)

func Mock(opts ...optFunc) *awsmocker.MockedEndpoint {

	options := &Options{
		Service: ecsTypes.Service{
			Status: aws.String("ACTIVE"),
		},
	}

	for _, optFunc := range opts {
		optFunc(options)
	}

	req := &awsmocker.MockedRequest{
		Service:       "ecs",
		Action:        "DescribeServices",
		MaxMatchCount: options.MaxCount,
	}

	if options.Name != "" {
		req.Matcher = func(rr *awsmocker.ReceivedRequest) bool {
			svcNameArnRaw := testutil.JmesSearchOrNil(rr.JsonPayload, "services[0]")
			if svcNameArnRaw == nil {
				return false
			}
			return options.Name == filepath.Base(svcNameArnRaw.(string))
		}
	}

	svc := options.Service

	return &awsmocker.MockedEndpoint{
		Request: req,
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				clusterArn := testutil.JmesPathSearch(rr.JsonPayload, "cluster").(string)
				serviceName := testutil.JmesPathSearch(rr.JsonPayload, "services[0]").(string)

				if options.Missing {
					return jsonify(map[string]any{
						"failures": []any{
							map[string]any{
								"arn":    fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, serviceName),
								"reason": "MISSING",
							},
						},
						"services": []any{},
					})
				}

				svc.ClusterArn = &clusterArn
				svc.ServiceName = &serviceName

				return jsonify(map[string]any{
					"failures": []any{},
					"services": []any{testutil.MockResponse_ECS_Service(svc)},
				})
			},
		},
	}
}

func jsonify(thing any) string {
	result, _ := util.Jsonify(thing)
	return result
}
