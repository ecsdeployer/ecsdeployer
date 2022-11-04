package taskdef

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// Builds log configuration, and returns a log config option and an additional container if needed (like firelens)
func ApplyLoggingConfiguration(input *pipelineInput) error {
	var logConf *ecsTypes.LogConfiguration
	var newCont *ecsTypes.ContainerDefinition
	var err error

	if input.Common.LoggingConfig == nil {
		logConf, newCont, err = loggingConfBuilderDefault(input)
	} else {
		logConf, newCont, err = loggingConfBuilderCustom(input)
	}

	if err != nil {
		return err
	}

	if logConf != nil {
		for i := range input.TaskDef.ContainerDefinitions {
			input.TaskDef.ContainerDefinitions[i].LogConfiguration = logConf
		}
	}

	if newCont != nil {

		for i := range input.TaskDef.ContainerDefinitions {
			if input.TaskDef.ContainerDefinitions[i].DependsOn == nil {
				input.TaskDef.ContainerDefinitions[i].DependsOn = make([]ecsTypes.ContainerDependency, 0, 1)
			}

			input.TaskDef.ContainerDefinitions[i].DependsOn = append(input.TaskDef.ContainerDefinitions[i].DependsOn, ecsTypes.ContainerDependency{
				Condition:     ecsTypes.ContainerConditionStart,
				ContainerName: newCont.Name,
			})
		}
		input.TaskDef.ContainerDefinitions = append(input.TaskDef.ContainerDefinitions, *newCont)
	}

	return nil
}

func loggingConfBuilderDefault(input *pipelineInput) (*ecsTypes.LogConfiguration, *ecsTypes.ContainerDefinition, error) {

	// if they disabled the default logging
	if input.Context.Project.Logging.IsDisabled() {
		// log.WithField("name", input.Common.Name).Debug("logging disabled for project")
		return nil, nil, nil
	}

	logConfig := input.Context.Project.Logging

	if !logConfig.FirelensConfig.IsDisabled() {
		return loggingConfBuilderDefaultFirelens(input)
	}

	return loggingConfBuilderDefaultAwslogs(input)
}
