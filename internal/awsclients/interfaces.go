package awsclients

import (
	"context"

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

type STSClienter interface {
	GetCallerIdentity(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type SSMClienter interface {
	GetParameter(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParameterHistory(context.Context, *ssm.GetParameterHistoryInput, ...func(*ssm.Options)) (*ssm.GetParameterHistoryOutput, error)
	GetParameters(context.Context, *ssm.GetParametersInput, ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
	GetParametersByPath(context.Context, *ssm.GetParametersByPathInput, ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}

type EC2Clienter interface {
	DescribeSecurityGroups(context.Context, *ec2.DescribeSecurityGroupsInput, ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeSubnets(context.Context, *ec2.DescribeSubnetsInput, ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
}

type ECSClienter interface {
	CreateService(context.Context, *ecs.CreateServiceInput, ...func(*ecs.Options)) (*ecs.CreateServiceOutput, error)
	DeleteService(context.Context, *ecs.DeleteServiceInput, ...func(*ecs.Options)) (*ecs.DeleteServiceOutput, error)
	DeregisterTaskDefinition(context.Context, *ecs.DeregisterTaskDefinitionInput, ...func(*ecs.Options)) (*ecs.DeregisterTaskDefinitionOutput, error)
	DescribeServices(context.Context, *ecs.DescribeServicesInput, ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	DescribeTaskDefinition(context.Context, *ecs.DescribeTaskDefinitionInput, ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
	DescribeTasks(context.Context, *ecs.DescribeTasksInput, ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ListServices(context.Context, *ecs.ListServicesInput, ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
	ListTaskDefinitionFamilies(context.Context, *ecs.ListTaskDefinitionFamiliesInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error)
	ListTaskDefinitions(context.Context, *ecs.ListTaskDefinitionsInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
	RegisterTaskDefinition(context.Context, *ecs.RegisterTaskDefinitionInput, ...func(*ecs.Options)) (*ecs.RegisterTaskDefinitionOutput, error)
	RunTask(context.Context, *ecs.RunTaskInput, ...func(*ecs.Options)) (*ecs.RunTaskOutput, error)
	UpdateService(context.Context, *ecs.UpdateServiceInput, ...func(*ecs.Options)) (*ecs.UpdateServiceOutput, error)
}

type ELBv2Clienter interface {
	DescribeTargetGroups(context.Context, *elbv2.DescribeTargetGroupsInput, ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error)
}

type EventsClienter interface {
	PutRule(context.Context, *events.PutRuleInput, ...func(*events.Options)) (*events.PutRuleOutput, error)
	PutTargets(context.Context, *events.PutTargetsInput, ...func(*events.Options)) (*events.PutTargetsOutput, error)
	ListTargetsByRule(context.Context, *events.ListTargetsByRuleInput, ...func(*events.Options)) (*events.ListTargetsByRuleOutput, error)
	RemoveTargets(context.Context, *events.RemoveTargetsInput, ...func(*events.Options)) (*events.RemoveTargetsOutput, error)
	DeleteRule(context.Context, *events.DeleteRuleInput, ...func(*events.Options)) (*events.DeleteRuleOutput, error)
}

type LogsClienter interface {
	CreateLogGroup(context.Context, *logs.CreateLogGroupInput, ...func(*logs.Options)) (*logs.CreateLogGroupOutput, error)
	DeleteRetentionPolicy(context.Context, *logs.DeleteRetentionPolicyInput, ...func(*logs.Options)) (*logs.DeleteRetentionPolicyOutput, error)
	DescribeLogGroups(context.Context, *logs.DescribeLogGroupsInput, ...func(*logs.Options)) (*logs.DescribeLogGroupsOutput, error)
	PutRetentionPolicy(context.Context, *logs.PutRetentionPolicyInput, ...func(*logs.Options)) (*logs.PutRetentionPolicyOutput, error)
}

type SchedulerClienter interface {
	CreateSchedule(context.Context, *scheduler.CreateScheduleInput, ...func(*scheduler.Options)) (*scheduler.CreateScheduleOutput, error)
	CreateScheduleGroup(context.Context, *scheduler.CreateScheduleGroupInput, ...func(*scheduler.Options)) (*scheduler.CreateScheduleGroupOutput, error)
	DeleteSchedule(context.Context, *scheduler.DeleteScheduleInput, ...func(*scheduler.Options)) (*scheduler.DeleteScheduleOutput, error)
	GetSchedule(context.Context, *scheduler.GetScheduleInput, ...func(*scheduler.Options)) (*scheduler.GetScheduleOutput, error)
	GetScheduleGroup(context.Context, *scheduler.GetScheduleGroupInput, ...func(*scheduler.Options)) (*scheduler.GetScheduleGroupOutput, error)
	ListSchedules(context.Context, *scheduler.ListSchedulesInput, ...func(*scheduler.Options)) (*scheduler.ListSchedulesOutput, error)
	UpdateSchedule(context.Context, *scheduler.UpdateScheduleInput, ...func(*scheduler.Options)) (*scheduler.UpdateScheduleOutput, error)
}

type TaggingClienter interface {
	GetResources(context.Context, *tagging.GetResourcesInput, ...func(*tagging.Options)) (*tagging.GetResourcesOutput, error)
}
