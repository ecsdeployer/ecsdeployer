package config

import (
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func (obj *AwsClientManager) STSClient() *sts.Client {
	return obj.stsClient
}

func (obj *AwsClientManager) SSMClient() *ssm.Client {
	return obj.ssmClient
}

func (obj *AwsClientManager) ECSClient() *ecs.Client {
	return obj.ecsClient
}

func (obj *AwsClientManager) EC2Client() *ec2.Client {
	return obj.ec2Client
}

func (obj *AwsClientManager) ELBv2Client() *elbv2.Client {
	return obj.elbv2Client
}

func (obj *AwsClientManager) LogsClient() *logs.Client {
	return obj.logsClient
}

func (obj *AwsClientManager) EventsClient() *events.Client {
	return obj.eventsClient
}

func (obj *AwsClientManager) TaggingClient() *tagging.Client {
	return obj.taggingClient
}

// func (obj *AwsClientManager) ELBClient() *elb.Client {
// 	return obj.elbClient
// }

// func (obj *AwsClientManager) IAMClient() *iam.Client {
// 	return obj.iamClient
// }
