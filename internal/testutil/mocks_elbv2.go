package testutil

import (
	"net/url"

	"github.com/webdestroya/awsmocker"
)

func Mock_ELBv2_DescribeTargetGroups_Single_Success(tgName string) *awsmocker.MockedEndpoint {
	return &awsmocker.MockedEndpoint{
		Request: &awsmocker.MockedRequest{
			Service: "elasticloadbalancing",
			Action:  "DescribeTargetGroups",
			Params: url.Values{
				"Names.member.1": []string{tgName},
			},
		},
		Response: MockResponse_ELBv2_DescribeTargetGroups_Single(tgName),
	}
}

func MockResponse_ELBv2_DescribeTargetGroups_Single(tgName string) *awsmocker.MockedResponse {

	return &awsmocker.MockedResponse{
		ContentType: awsmocker.ContentTypeXML,
		Body: TemplateApply(mock_tpl_DescribeTargetGroupsResponse, map[string]interface{}{
			"TargetGroupNames": []string{tgName},
			"AccountId":        awsmocker.DefaultAccountId,
		}),
	}
}

const (
	mock_tpl_DescribeTargetGroupsResponse = `<DescribeTargetGroupsResponse xmlns="http://elasticloadbalancing.amazonaws.com/doc/2015-12-01/">
  <DescribeTargetGroupsResult> 
    <TargetGroups>
			{{ range .TargetGroupNames }}
      <member> 
        <TargetGroupArn>arn:aws:elasticloadbalancing:us-east-1:{{$.AccountId}}:targetgroup/{{.}}/73e2d6bc24d8a067</TargetGroupArn> 
        <HealthCheckTimeoutSeconds>5</HealthCheckTimeoutSeconds> 
        <HealthCheckPort>traffic-port</HealthCheckPort> 
        <Matcher> 
          <HttpCode>200</HttpCode> 
        </Matcher> 
        <TargetGroupName>{{.}}</TargetGroupName> 
        <HealthCheckProtocol>HTTP</HealthCheckProtocol> 
        <HealthCheckPath>/</HealthCheckPath> 
        <Protocol>HTTP</Protocol> 
        <Port>80</Port> 
        <VpcId>vpc-3ac0fb5f</VpcId> 
        <HealthyThresholdCount>5</HealthyThresholdCount> 
        <HealthCheckIntervalSeconds>30</HealthCheckIntervalSeconds> 
        <LoadBalancerArns> 
          <member>arn:aws:elasticloadbalancing:us-east-1:{{$.AccountId}}:loadbalancer/app/my-load-balancer/50dc6c495c0c9188</member> 
        </LoadBalancerArns> 
        <UnhealthyThresholdCount>2</UnhealthyThresholdCount> 
      </member>
			{{ end }}
    </TargetGroups> 
  </DescribeTargetGroupsResult> 
  <ResponseMetadata> 
    <RequestId>70092c0e-f3a9-11e5-ae48-cff02092876b</RequestId> 
  </ResponseMetadata> 
</DescribeTargetGroupsResponse>`
)
