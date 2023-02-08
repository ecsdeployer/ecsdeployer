package awsclients

import (
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func STSClient() *sts.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return stsClient
}

func SSMClient() *ssm.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ssmClient
}

func ECSClient() *ecs.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ecsClient
}

func EC2Client() *ec2.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ec2Client
}

func ELBv2Client() *elbv2.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return elbv2Client
}

func LogsClient() *logs.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return logsClient
}

func EventsClient() *events.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return eventsClient
}

func TaggingClient() *tagging.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return taggingClient
}

func SchedulerClient() *scheduler.Client {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return schedulerClient
}
