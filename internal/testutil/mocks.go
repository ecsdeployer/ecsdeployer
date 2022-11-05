package testutil

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/webdestroya/awsmocker"
)

// This is just a basic mock server to get the account ID and region
func MockSimpleStsProxy(t *testing.T) func() {
	// awsmocker.GlobalDebugMode = true
	closeFunc, _, _ := awsmocker.StartMockServer(&awsmocker.MockerOptions{
		T: t,
	})
	return closeFunc
}

// func Mock_EC2_DescribeSubnets(params map[string]string, )

func Mock_ELBv2_DescribeTargetGroups_Single_Success(tgName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "elasticloadbalancing",
			Action:  "DescribeTargetGroups",
			Params: url.Values{
				"Names.member.1": []string{tgName},
			},
		},
		Response: &awsmocker.MockedResponse{
			ContentType: awsmocker.ContentTypeXML,
			Body:        fmt.Sprintf(mock_DescribeTargetGroupsResponse_SingleResult, tgName, tgName),
		},
	}
}

const (
	mock_DescribeTargetGroupsResponse_SingleResult = `<DescribeTargetGroupsResponse xmlns="http://elasticloadbalancing.amazonaws.com/doc/2015-12-01/">
  <DescribeTargetGroupsResult> 
    <TargetGroups> 
      <member> 
        <TargetGroupArn>arn:aws:elasticloadbalancing:us-east-1:555555555555:targetgroup/%s/73e2d6bc24d8a067</TargetGroupArn> 
        <HealthCheckTimeoutSeconds>5</HealthCheckTimeoutSeconds> 
        <HealthCheckPort>traffic-port</HealthCheckPort> 
        <Matcher> 
          <HttpCode>200</HttpCode> 
        </Matcher> 
        <TargetGroupName>%s</TargetGroupName> 
        <HealthCheckProtocol>HTTP</HealthCheckProtocol> 
        <HealthCheckPath>/</HealthCheckPath> 
        <Protocol>HTTP</Protocol> 
        <Port>80</Port> 
        <VpcId>vpc-3ac0fb5f</VpcId> 
        <HealthyThresholdCount>5</HealthyThresholdCount> 
        <HealthCheckIntervalSeconds>30</HealthCheckIntervalSeconds> 
        <LoadBalancerArns> 
          <member>arn:aws:elasticloadbalancing:us-east-1:555555555555:loadbalancer/app/my-load-balancer/50dc6c495c0c9188</member> 
        </LoadBalancerArns> 
        <UnhealthyThresholdCount>2</UnhealthyThresholdCount> 
      </member> 
    </TargetGroups> 
  </DescribeTargetGroupsResult> 
  <ResponseMetadata> 
    <RequestId>70092c0e-f3a9-11e5-ae48-cff02092876b</RequestId> 
  </ResponseMetadata> 
</DescribeTargetGroupsResponse>`
)
