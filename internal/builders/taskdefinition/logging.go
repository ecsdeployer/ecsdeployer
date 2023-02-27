package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyContainerLogging(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	if b.project.Logging.IsDisabled() {
		// logging was disabled at the project level
		return nil
	}

	common := thing.GetCommonContainerAttrs()

	loggingConf := util.Coalesce(common.LoggingConfig, b.entity.GetCommonContainerAttrs().LoggingConfig)

	if loggingConf != nil {
		// this thing has fancy custom logging
		return b.buildContainerLogging(cdef, loggingConf)
	}

	// generic logging setup
	return b.applyContainerLoggingDefault(cdef, thing)
}

func (b *Builder) applyContainerLoggingDefault(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {
	logConfig := b.project.Logging

	if !logConfig.Custom.IsDisabled() {
		return b.applyContainerLoggingCustom(cdef, thing)
	}

	if !logConfig.FirelensConfig.IsDisabled() {
		return b.applyContainerLoggingFirelens(cdef, thing)
	}

	return b.applyContainerLoggingAwsLogs(cdef, thing)
}

func (b *Builder) buildContainerLogging(cdef *ecsTypes.ContainerDefinition, logConfig *config.TaskLoggingConfig) error {
	if logConfig == nil {
		return nil
	}

	if logConfig.IsDisabled() || logConfig.Driver == nil {
		return nil
	}

	conf := &ecsTypes.LogConfiguration{
		LogDriver:     ecsTypes.LogDriver(*logConfig.Driver),
		Options:       make(map[string]string),
		SecretOptions: make([]ecsTypes.Secret, 0),
	}

	tpl := b.containerTpl(cdef)

	for lk, lv := range logConfig.Options.Filter() {
		if lv.IsSSM() {

			conf.SecretOptions = append(conf.SecretOptions, ecsTypes.Secret{
				Name:      aws.String(lk),
				ValueFrom: aws.String(util.Must(lv.GetValue(nil))),
			})
			continue
		}
		val, err := lv.GetValue(tpl)
		if err != nil {
			return err
		}
		conf.Options[lk] = val
	}

	if len(conf.Options) == 0 {
		conf.Options = nil
	}

	if len(conf.SecretOptions) == 0 {
		conf.SecretOptions = nil
	}

	cdef.LogConfiguration = conf

	return nil
}
