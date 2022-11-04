package config_test

import (
	"context"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestContextClients_Smoke(t *testing.T) {
	mgr := config.NewAwsClientManager(context.TODO())

	_ = mgr.STSClient()
	_ = mgr.SSMClient()
	_ = mgr.ECSClient()
	_ = mgr.EC2Client()
	_ = mgr.ELBv2Client()
	_ = mgr.LogsClient()
	_ = mgr.EventsClient()
	_ = mgr.TaggingClient()

}
