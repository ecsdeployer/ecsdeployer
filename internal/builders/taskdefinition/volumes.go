package taskdefinition

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyVolumes() error {

	volumeList := config.MergeVolumeLists(b.taskDefaults.Volumes, b.commonTask.Volumes)

	if len(volumeList) == 0 {
		return nil
	}

	mounts := util.DeepFindInStruct[ecsTypes.MountPoint](b.taskDef)

	if len(mounts) == 0 {
		return nil
	}

	usedVolumes := make(config.VolumeList)

	for _, mount := range mounts {
		volName := *mount.SourceVolume

		vol, ok := volumeList[volName]
		if !ok {
			return fmt.Errorf("Missing volume declaration for '%s'", volName)
		}

		usedVolumes[volName] = vol
	}

	b.taskDef.Volumes = usedVolumes.ToAws()

	return nil
}
