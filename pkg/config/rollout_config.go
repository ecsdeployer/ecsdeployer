package config

import (
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type RolloutConfig struct {
	Minimum *int32 `yaml:"min" json:"min" jsonschema:"required" jsonschema_extras:"minimum=0"`
	Maximum *int32 `yaml:"max" json:"max" jsonschema:"required" jsonschema_extras:"minimum=0"`
}

func (dc *RolloutConfig) GetAwsConfig() *ecsTypes.DeploymentConfiguration {
	return &ecsTypes.DeploymentConfiguration{
		MinimumHealthyPercent: dc.Minimum,
		MaximumPercent:        dc.Maximum,
	}
}

func (obj *RolloutConfig) MaximumPercent() float64 {
	return float64(*obj.Maximum) / 100.0
}

func (obj *RolloutConfig) MinimumPercent() float64 {
	return float64(*obj.Minimum) / 100.0
}

func (obj *RolloutConfig) GetMinMaxCount(count int32) (int32, int32) {
	min := int32(math.Ceil(obj.MinimumPercent() * float64(count)))
	max := int32(math.Floor(obj.MaximumPercent() * float64(count)))

	return min, max
}

func (a *RolloutConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tRolloutConfig RolloutConfig
	var obj = tRolloutConfig{}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	*a = RolloutConfig(obj)

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (def *RolloutConfig) Validate() error {
	if def.Maximum == nil || def.Minimum == nil {
		return NewValidationError("you must define both 'min' and 'max' for rollout configuration")
	}

	if *def.Minimum >= *def.Maximum {
		return NewValidationError("RolloutConfiguration: maximum must be higher than minimum")
	}

	if *def.Maximum < 0 || *def.Minimum < 0 {
		return NewValidationError("RolloutConfig: Min/Max cannot be below zero")
	}

	return nil
}

func (obj *RolloutConfig) ValidateWithDesiredCount(count int32) error {
	if err := obj.Validate(); err != nil {
		return err
	}

	// zero count doesnt matter
	if count == 0 {
		return nil
	}

	minHealthyCount, maxHealthyCount := obj.GetMinMaxCount(count)

	if minHealthyCount >= maxHealthyCount {
		return NewValidationError("this is an impossible rollout config. Desired=%d, Min=%d, Max=%d. This would give minCount=%d and maxCount=%d",
			count,
			*obj.Minimum,
			*obj.Maximum,
			minHealthyCount,
			maxHealthyCount,
		)
	}

	return nil
}

func NewDeploymentConfigFromService(svc *Service) *RolloutConfig {
	desired := svc.DesiredCount

	if desired == 0 {
		return &RolloutConfig{Minimum: aws.Int32(0), Maximum: aws.Int32(100)}
	}

	if svc.IsLoadBalanced() {
		if desired > 4 {
			return &RolloutConfig{Minimum: aws.Int32(100), Maximum: aws.Int32(150)}
		} else {
			return &RolloutConfig{Minimum: aws.Int32(100), Maximum: aws.Int32(200)}
		}
	}

	if desired == 1 {
		return &RolloutConfig{Minimum: aws.Int32(100), Maximum: aws.Int32(200)}
	}

	return &RolloutConfig{Minimum: aws.Int32(50), Maximum: aws.Int32(100)}
}
