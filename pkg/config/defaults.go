package config

import "regexp"

const (
	defaultConsoleEnabled = false
	defaultConsolePort    = 8722

	defaultTaskCpu    = 1024
	defaultTaskMemory = "2x"

	defaultKeepInSync = true

	defaultLogRetention = 180
)

var DefaultDeploymentEnvVars = map[string]string{
	"ECSDEPLOYER_PROJECT":     "{{ .ProjectName }}",
	"ECSDEPLOYER_TASK_NAME":   "{{ .Name }}",
	"ECSDEPLOYER_STAGE":       "{{ .Stage }}",
	"ECSDEPLOYER_DEPLOYED_AT": "{{ .Date }}",
	"ECSDEPLOYER_APP_VERSION": "{{ .Version }}",
	"ECSDEPLOYER_IMAGE_TAG":   "{{ .ImageTag }}",
}

// Regex to validate names of things
const shortCodeNameRegexStr = "^[a-z][-_a-z0-9]*$"

var shortCodeNameRegex = regexp.MustCompile(shortCodeNameRegexStr)

// func DefaultDeploymentEnvVars() map[string]string {
// 	return defaultDeploymentEnvVars
// }
