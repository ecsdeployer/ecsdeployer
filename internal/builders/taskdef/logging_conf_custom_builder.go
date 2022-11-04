package taskdef

import (
	"errors"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func loggingConfBuilderCustom(input *pipelineInput) (*ecsTypes.LogConfiguration, *ecsTypes.ContainerDefinition, error) {
	logConfig := input.Common.LoggingConfig

	if logConfig.IsDisabled() {
		// log.WithField("name", input.Common.Name).Debug("logging disabled for task")
		return nil, nil, nil
	}

	// TODO: FINISH THIS
	return nil, nil, errors.New("Custom logging is not yet supported")

	// return nil, nil, nil
}
