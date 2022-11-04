package helpers

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func GetTemplatedPrefix(ctx *config.Context, tplStr string) (string, error) {
	tpl := tmpl.New(ctx).WithExtraFields(tmpl.Fields{
		"Name":        "",
		"TaskName":    "",
		"ServiceName": "",
	})

	tplValue, err := tpl.Apply(tplStr)
	if err != nil {
		return "", err
	}

	return tplValue, nil
}

func GetDefaultTaskTemplateFields(ctx *config.Context, common *config.CommonTaskAttrs) (tmpl.Fields, error) {

	project := ctx.Project

	fields := tmpl.Fields{
		"TaskName": common.Name,
		"Name":     common.Name,
	}

	arch := util.Coalesce(common.Architecture, project.TaskDefaults.Architecture)
	if arch != nil {
		fields["Arch"] = string(*arch)
	}

	return fields, nil
}
