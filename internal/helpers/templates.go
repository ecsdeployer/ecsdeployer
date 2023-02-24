package helpers

import (
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

const templatePrefixSentinel = "@ECSD_REPLACE@"

// This should get the first part of a template string before any dynamic fields exist.
// so it can include project name, stage, cluster, but it is cut at the point where it would
// reference a service/name/task
func GetTemplatedPrefix(ctx *config.Context, tplStr string) (string, error) {
	tpl := tmpl.New(ctx).WithExtraFields(tmpl.Fields{
		"Name":        templatePrefixSentinel,
		"TaskName":    templatePrefixSentinel,
		"ServiceName": templatePrefixSentinel,
	})

	tplValue, err := tpl.Apply(tplStr)
	if err != nil {
		return "", err
	}

	prefix, _, _ := strings.Cut(tplValue, templatePrefixSentinel)

	return prefix, nil
}

func GetDefaultTaskTemplateFields(ctx *config.Context, common *config.CommonTaskAttrs) (tmpl.Fields, error) {

	project := ctx.Project

	fields := tmpl.Fields{
		"TaskName": common.Name,
		"Name":     common.Name,
	}

	arch := util.Coalesce(common.Architecture, project.TaskDefaults.Architecture)
	if arch != nil {
		fields["Arch"] = arch.String()
	}

	return fields, nil
}
