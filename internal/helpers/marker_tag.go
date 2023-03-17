package helpers

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

// Returns (key, value, err)
func GetMarkerTag(ctx *config.Context) (string, string, error) {
	tpl := tmpl.New(ctx)

	keyVal, err := tpl.Apply(*ctx.Project.Templates.MarkerTagKey)
	if err != nil {
		return "", "", err
	}

	valVal, err := tpl.Apply(*ctx.Project.Templates.MarkerTagValue)
	if err != nil {
		return "", "", err
	}

	return keyVal, valVal, nil
}
