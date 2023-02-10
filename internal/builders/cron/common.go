package cron

import "ecsdeployer.com/ecsdeployer/pkg/config"

type cronContainerOverride struct {
	Name    string              `json:"name"`
	Command config.ShellCommand `json:"command,omitempty"`
}

type cronInputObj struct {
	ContainerOverrides []cronContainerOverride `json:"containerOverrides,omitempty"`
}
