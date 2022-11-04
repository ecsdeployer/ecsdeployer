package helpers

import (
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/util"
	awsarn "github.com/aws/aws-sdk-go-v2/aws/arn"
)

func GetECSServiceNameFromArn(str string) string {

	if !IsArnForService(str, "ecs") {
		return ""
	}
	arn := util.Must(awsarn.Parse(str))

	res := arn.Resource

	if !strings.HasPrefix(res, "service/") {
		// wrong kinda thing
		return ""
	}

	trimmedRes := strings.TrimPrefix(res, "service/")

	before, after, hasSlash := strings.Cut(trimmedRes, "/")

	if hasSlash {
		// "service/" "cluster/serviceName"
		return after
	} else {
		// "service/" "serviceName"
		return before
	}

}

func GetECSClusterNameFromArn(str string) string {

	if !IsArnForService(str, "ecs") {
		return ""
	}
	arn := util.Must(awsarn.Parse(str))

	res := arn.Resource

	// it's a cluster ARN
	if strings.HasPrefix(res, "cluster/") {
		return strings.TrimPrefix(res, "cluster/")
	}

	parts := strings.Split(res, "/")

	if len(parts) < 3 {
		return ""
	}

	switch parts[0] {
	case "service", "task-set", "container-instance", "task":
		return parts[1]
	default:
		return ""
	}
}

func GetTaskDefFamilyFromArn(str string) string {
	if !IsArnForService(str, "ecs") {
		return ""
	}
	arn := util.Must(awsarn.Parse(str))

	res := arn.Resource

	if !strings.HasPrefix(res, "task-definition/") {
		return ""
	}

	before, _, _ := strings.Cut(strings.TrimPrefix(res, "task-definition/"), ":")
	return before
}

func IsArnForService(str string, serviceCode string) bool {
	if !awsarn.IsARN(str) {
		return false
	}

	arn, err := awsarn.Parse(str)
	if err != nil {
		return false
	}

	if arn.Service != serviceCode {
		return false
	}

	return true
}

// return = (ruleName, eventBus)
func GetEventRuleNameAndBusFromArn(str string) (string, string) {
	if !IsArnForService(str, "events") {
		return "", ""
	}
	arn := util.Must(awsarn.Parse(str))

	res := arn.Resource

	if !strings.HasPrefix(res, "rule/") {
		return "", ""
	}

	part1, part2, hasEventBus := strings.Cut(strings.TrimPrefix(res, "rule/"), "/")
	if hasEventBus {
		return part2, part1
	} else {
		return part1, ""
	}
}
