package cronschedule

import "ecsdeployer.com/ecsdeployer/pkg/config"

type cronOverrideKeyPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type cronContainerOverride struct {
	Name        string                `json:"name"`
	Command     config.ShellCommand   `json:"command,omitempty"`
	Environment []cronOverrideKeyPair `json:"environment,omitempty"`
}

type cronInputObj struct {
	ContainerOverrides []cronContainerOverride `json:"containerOverrides,omitempty"`
}
