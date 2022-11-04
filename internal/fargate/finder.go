package fargate

func FindFargateBestFit(cpu int32, memory int32) FargateResource {

	for _, res := range DefaultFargateResources {
		if res.Fits(cpu, memory) {
			return res
		}
	}

	return DefaultFargateResources[len(DefaultFargateResources)-1]
}

// Are these requirements excessively big?
// the size they requested is too big even for the largest known instance
func ExceedsLargest(cpu int32, memory int32) bool {
	biggest := DefaultFargateResources[len(DefaultFargateResources)-1]

	return !biggest.Fits(cpu, memory)
}

// Returns the best fargate size, or makes a new fargate size and trusts the user
func FindFargateBestFitOrTrust(cpu int32, memory int32) FargateResource {

	if ExceedsLargest(cpu, memory) {
		return FargateResource{
			Cpu:    cpu,
			Memory: memory,
		}
	}

	return FindFargateBestFit(cpu, memory)
}
