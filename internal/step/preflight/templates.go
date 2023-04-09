package preflight

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkTemplates struct{}

func (checkTemplates) String() string {
	return "templates"
}

func (checkTemplates) Check(ctx *config.Context) error {
	tpl := tmpl.New(ctx).WithExtraFields(tmpl.Fields{
		"Name":      "THING",
		"Container": "THING",
		"Arch":      "amd64",
	})

	for _, val := range util.DeepFindInStruct[string](ctx.Project.Templates) {
		_, err := tpl.Apply(*val)
		if err != nil {
			return err
		}
	}

	return nil
}
