package awsclients

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
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

var (
	initMutex sync.RWMutex

	awsConfig aws.Config

	ecsClient       ECSClienter
	stsClient       STSClienter
	ssmClient       SSMClienter
	ec2Client       EC2Clienter
	elbv2Client     ELBv2Clienter
	eventsClient    EventsClienter
	logsClient      LogsClienter
	taggingClient   TaggingClienter
	schedulerClient SchedulerClienter
)

func init() {

	httpClient := awsHttp.NewBuildableClient().WithTransportOptions(func(t *http.Transport) {
		t.ResponseHeaderTimeout = 30 * time.Second
	})

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		// panic(fmt.Sprintf("failed loading config, %v", err))
		// couldnt load default? maybe they have a better one
		return
	}

	// use default config initially
	SetupWithConfig(cfg)
}

func SetupWithConfig(cfg aws.Config) {
	initMutex.Lock()
	defer initMutex.Unlock()

	awsConfig = cfg

	stsClient = sts.NewFromConfig(cfg)
	ecsClient = ecs.NewFromConfig(cfg)
	ssmClient = ssm.NewFromConfig(cfg)
	taggingClient = tagging.NewFromConfig(cfg)
	ec2Client = ec2.NewFromConfig(cfg)
	elbv2Client = elbv2.NewFromConfig(cfg)
	eventsClient = events.NewFromConfig(cfg)
	logsClient = logs.NewFromConfig(cfg)
	schedulerClient = scheduler.NewFromConfig(cfg)
}

func AwsConfig() aws.Config {
	return awsConfig
}
