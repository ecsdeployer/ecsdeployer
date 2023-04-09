package config

import (
	"sync"

	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type ContextCache struct {
	SSMSecretsCached bool
	SSMSecrets       map[string]EnvVar

	LogGroupsCached bool
	LogGroups       map[string]logTypes.LogGroup

	registeredTaskDefArns []string

	// TODO: secretsmanager secrets?
	// TODO: task families?
	// TODO: services?
	// TODO: cron rules?
	// TODO: cron targets?

	Meta map[string]interface{}

	muTaskDefs sync.Mutex
}

func (cc *ContextCache) AddTaskDefinition(arn string) {
	cc.muTaskDefs.Lock()
	defer cc.muTaskDefs.Unlock()

	cc.registeredTaskDefArns = append(cc.registeredTaskDefArns, arn)
}

func (cc *ContextCache) TaskDefinitions() []string {
	return cc.registeredTaskDefArns
}

func newContextCache() *ContextCache {
	return &ContextCache{
		SSMSecrets:            make(map[string]EnvVar),
		LogGroups:             make(map[string]logTypes.LogGroup, 0),
		registeredTaskDefArns: make([]string, 0),
	}
}
