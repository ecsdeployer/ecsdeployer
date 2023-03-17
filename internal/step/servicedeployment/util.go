package servicedeployment

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func getServiceName(ctx *config.Context, service *config.Service) (string, error) {
	serviceName, err := tmpl.New(ctx).WithExtraFields(service.TemplateFields()).Apply(*ctx.Project.Templates.ServiceName)
	if err != nil {
		return "", err
	}

	return serviceName, nil
}
