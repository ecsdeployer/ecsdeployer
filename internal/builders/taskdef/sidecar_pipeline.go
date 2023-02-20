package taskdef

import (
	"errors"
)

func SidecarPipeline(input *pipelineInput) error {

	// TODO: add sidecar containers
	common := input.Common

	if len(common.Sidecars) == 0 {
		return nil
	}
	return errors.New("SIDECAR PIPELINE NOT YET FINISHED")
}
