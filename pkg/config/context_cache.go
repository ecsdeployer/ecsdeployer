package config

import (
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type ContextCache struct {
	SSMSecrets map[string]EnvVar

	LogGroups []logTypes.LogGroup

	RegisteredTaskDefArns []string

	// TODO: secretsmanager secrets?
	// TODO: task families?
	// TODO: services?
	// TODO: cron rules?
	// TODO: cron targets?

	Meta map[string]interface{}
}

func newContextCache() *ContextCache {
	return &ContextCache{
		SSMSecrets:            make(map[string]EnvVar),
		LogGroups:             make([]logTypes.LogGroup, 0),
		RegisteredTaskDefArns: make([]string, 0),
	}
}
