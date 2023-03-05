package taskdefinition

import ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"

func hasContainerDependency(cdef *ecsTypes.ContainerDefinition, dependsOn string) bool {
	depList := cdef.DependsOn

	if len(depList) == 0 {
		return false
	}

	for _, entry := range depList {
		if *entry.ContainerName == dependsOn {
			return true
		}
	}

	return false
}

// dependsOn can be string, *string, containerDef
func addContainerDependency(cdef *ecsTypes.ContainerDefinition, dependsOn any, condition ecsTypes.ContainerCondition) {

	depOnName := ""

	switch val := dependsOn.(type) {
	case string:
		depOnName = val
	case *string:
		depOnName = *val
	case *ecsTypes.ContainerDefinition:
		depOnName = *val.Name
	case ecsTypes.ContainerDefinition:
		depOnName = *val.Name
	default:
		panic("wrong type for containerdependency.dependsOn")
	}

	if cdef.DependsOn == nil {
		cdef.DependsOn = make([]ecsTypes.ContainerDependency, 0)
	}

	if hasContainerDependency(cdef, depOnName) {
		return
	}

	cdef.DependsOn = append(cdef.DependsOn, ecsTypes.ContainerDependency{
		ContainerName: &depOnName,
		Condition:     condition,
	})
}
