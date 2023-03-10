package helpers

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func ResolveImageUri(ctx *config.Context, img *config.ImageUri) (string, error) {

	if img.IsResolved() {
		return img.Value(), nil
	}

	newUri, err := tmpl.New(ctx).Apply(img.Value())
	if err != nil {
		return "", fmt.Errorf("failed to resolve Image URI: %w", err)
	}

	// log.WithField("image", newUri).Debug("Resolved Image")
	img.SetResolved(newUri)

	return newUri, nil
}
