package config

import (
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type ContextCache struct {
	SSMSecrets map[string]EnvVar

	LogGroups []logTypes.LogGroup

	// TODO: secretsmanager secrets?
	// TODO: task families?
	// TODO: services?
	// TODO: cron rules?
	// TODO: cron targets?

	Meta map[string]interface{}
}
