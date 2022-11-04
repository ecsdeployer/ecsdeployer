package helpers

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/log"
)

func ResolveImageUri(ctx *config.Context, img *config.ImageUri) (string, error) {

	if img.IsResolved() {
		return img.Value(), nil
	}

	newUri, err := tmpl.New(ctx).Apply(img.Value())
	if err != nil {
		return "", err
	}

	log.WithField("image", newUri).Debugf("Resolved Image")
	img.SetResolved(newUri)

	return newUri, nil
}
