package describeservicemock

import (
	"fmt"

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

	jmesMatches := map[string]any{}
	if options.Name != "" {
		jmesMatches["services[0]"] = options.Name
	}

	req := &awsmocker.MockedRequest{
		Service:       "ecs",
		Action:        "DescribeServices",
		MaxMatchCount: options.MaxCount,
	}

	if len(jmesMatches) > 0 {
		req.Matcher = testutil.JmesRequestMatcher(jmesMatches)
	}

	svc := options.Service

	return &awsmocker.MockedEndpoint{
		Request: req,
		Response: &awsmocker.MockedResponse{
			Body: func(rr *awsmocker.ReceivedRequest) string {

				clusterArn := testutil.JmesPathSearch(rr.JsonPayload, "cluster").(string)
				serviceName := testutil.JmesPathSearch(rr.JsonPayload, "services[0]").(string)

				if options.Missing {
					return jsonify(map[string]interface{}{
						"failures": []interface{}{
							map[string]interface{}{
								"arn":    fmt.Sprintf("arn:aws:ecs:%s:%s:service/%s", awsmocker.DefaultRegion, awsmocker.DefaultAccountId, serviceName),
								"reason": "MISSING",
							},
						},
						"services": []interface{}{},
					})
				}

				svc.ClusterArn = &clusterArn
				svc.ServiceName = &serviceName

				return jsonify(map[string]interface{}{
					"failures": []interface{}{},
					"services": []interface{}{testutil.MockResponse_ECS_Service(svc)},
				})
			},
		},
	}
}

func jsonify(thing any) string {
	result, _ := util.Jsonify(thing)
	return result
}
