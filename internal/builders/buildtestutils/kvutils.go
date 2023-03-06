package buildtestutils

import (
	"fmt"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type kvToMapFunc[T any] func(T) (string, string)

func KVListToMap[T any](kvList []T, mapFunc kvToMapFunc[T]) map[string]string {
	newMap := make(map[string]string, len(kvList))

	for _, entry := range kvList {
		k, v := mapFunc(entry)
		newMap[k] = v
	}

	return newMap
}

func KVListToMap_KVP(val ecsTypes.KeyValuePair) (string, string) {
	return *val.Name, *val.Value
}

func KVListToMap_Secret(val ecsTypes.Secret) (string, string) {
	return *val.Name, *val.ValueFrom
}

func KVListToMap_Depends(val ecsTypes.ContainerDependency) (string, string) {
	return *val.ContainerName, string(val.Condition)
}

type kvToSliceFunc[T any, K any] func(T) K

func KVListToSlice[T any, K any](kvList []T, sliceFunc kvToSliceFunc[T, K]) []K {
	newMap := make([]K, 0, len(kvList))

	for _, entry := range kvList {
		newMap = append(newMap, sliceFunc(entry))
	}

	return newMap
}

func KVListToSlice_PortMaps(val ecsTypes.PortMapping) string {
	return fmt.Sprintf("%d/%s", *val.ContainerPort, val.Protocol)
}
