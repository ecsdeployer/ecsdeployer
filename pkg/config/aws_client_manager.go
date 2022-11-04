package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AwsClientManager struct {
	// awsAccountId string
	awsConfig aws.Config

	ecsClient     *ecs.Client
	stsClient     *sts.Client
	ssmClient     *ssm.Client
	ec2Client     *ec2.Client
	elbv2Client   *elbv2.Client
	eventsClient  *events.Client
	logsClient    *logs.Client
	taggingClient *tagging.Client
	// iamClient     *iam.Client
	// elbClient     *elb.Client
}

func NewAwsClientManager(ctx context.Context) *AwsClientManager {

	// if fakeaws.IsTestMode() {
	// 	return NewAwsClientManagerFromConfig(fakeaws.GetAwsConfig())
	// }

	cfg, err := config.LoadDefaultConfig(ctx, func(lo *config.LoadOptions) error {
		// lo.ClientLogMode = &clientLogMode
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	return NewAwsClientManagerFromConfig(cfg)
}

func NewAwsClientManagerFromConfig(cfg aws.Config) *AwsClientManager {

	// cfg.ClientLogMode = aws.LogSigning | aws.LogRequest | aws.LogRequestWithBody | aws.LogRequestEventMessage

	return &AwsClientManager{
		awsConfig:     cfg,
		stsClient:     sts.NewFromConfig(cfg),
		ecsClient:     ecs.NewFromConfig(cfg),
		ssmClient:     ssm.NewFromConfig(cfg),
		taggingClient: tagging.NewFromConfig(cfg),
		ec2Client:     ec2.NewFromConfig(cfg),
		elbv2Client:   elbv2.NewFromConfig(cfg),
		eventsClient:  events.NewFromConfig(cfg),
		logsClient:    logs.NewFromConfig(cfg),
		// iamClient:     iam.NewFromConfig(cfg),
		// elbClient:     elb.NewFromConfig(cfg),
	}
}

func (obj *AwsClientManager) AwsConfig() aws.Config {
	return obj.awsConfig
}
