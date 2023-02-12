package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

const (
	capProviderSpot     = "FARGATE_SPOT"
	capProviderOnDemand = "FARGATE"
)

type SpotOverrides struct {
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema:"default=false,description=Enable Fargate Spot"`

	MinimumOnDemand        *int32 `yaml:"minimum_ondemand,omitempty" json:"minimum_ondemand,omitempty"`
	MinimumOnDemandPercent *int32 `yaml:"minimum_ondemand_percent,omitempty" json:"minimum_ondemand_percent,omitempty"`
}

func NewSpotOnDemand() *SpotOverrides {
	return &SpotOverrides{Enabled: false}
}

func (obj *SpotOverrides) IsDisabled() bool {
	return !obj.Enabled
}

func (obj *SpotOverrides) ExportCapacityStrategyEventBridge() []eventTypes.CapacityProviderStrategyItem {
	ecs := obj.ExportCapacityStrategy()
	eb := make([]eventTypes.CapacityProviderStrategyItem, len(ecs))

	for i, elm := range ecs {
		eb[i] = eventTypes.CapacityProviderStrategyItem{
			CapacityProvider: elm.CapacityProvider,
			Base:             elm.Base,
			Weight:           elm.Weight,
		}
	}

	return eb
}

func (obj *SpotOverrides) ExportCapacityStrategyScheduler() []schedulerTypes.CapacityProviderStrategyItem {
	ecs := obj.ExportCapacityStrategy()
	eb := make([]schedulerTypes.CapacityProviderStrategyItem, len(ecs))

	for i, elm := range ecs {
		eb[i] = schedulerTypes.CapacityProviderStrategyItem{
			CapacityProvider: elm.CapacityProvider,
			Base:             elm.Base,
			Weight:           elm.Weight,
		}
	}

	return eb
}

func (obj *SpotOverrides) ExportCapacityStrategy() []ecsTypes.CapacityProviderStrategyItem {

	if obj.IsDisabled() {
		return []ecsTypes.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String(capProviderOnDemand),
				Weight:           1,
				Base:             0,
			},
		}
	}

	if !obj.WantsOnDemand() {
		return []ecsTypes.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String(capProviderSpot),
				Weight:           1,
				Base:             0,
			},
		}
	}

	onDemandCap := ecsTypes.CapacityProviderStrategyItem{
		CapacityProvider: aws.String(capProviderOnDemand),
		Weight:           1,
		Base:             0,
	}

	if obj.MinimumOnDemand != nil {
		onDemandCap.Base = *obj.MinimumOnDemand
	}

	if obj.MinimumOnDemandPercent != nil {
		onDemandCap.Weight = *obj.MinimumOnDemandPercent
	}

	return []ecsTypes.CapacityProviderStrategyItem{
		{
			CapacityProvider: aws.String(capProviderSpot),
			Weight:           100,
			Base:             0,
		},
		onDemandCap,
	}
}

func (obj *SpotOverrides) WantsOnDemand() bool {

	if obj.IsDisabled() {
		return true
	}

	if obj.MinimumOnDemand != nil && *obj.MinimumOnDemand > 0 {
		return true
	}

	if obj.MinimumOnDemandPercent != nil && *obj.MinimumOnDemandPercent > 0 {
		return true
	}

	return false
}

func (obj *SpotOverrides) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var val bool
	if err := unmarshal(&val); err != nil {
		type t SpotOverrides // prevent recursive overflow
		var defo = t{}
		if err := unmarshal(&defo); err != nil {
			return err
		}

		*obj = SpotOverrides(defo)

	} else {
		*obj = SpotOverrides{
			Enabled: val,
		}
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *SpotOverrides) Validate() error {

	// if obj.MinimumOnDemand != nil && obj.MinimumOnDemandPercent != nil {
	// 	return NewValidationError("You cannot set minimum_ondemand and minimum_ondemand_percent at the same time.")
	// }

	return nil
}

func (obj *SpotOverrides) ApplyDefaults() {

}

func (SpotOverrides) JSONSchema() *jsonschema.Schema {

	properties := orderedmap.New()
	properties.Set("enabled", &jsonschema.Schema{
		Type:        "boolean",
		Default:     false,
		Description: "Enable spot containers",
	})

	properties.Set("minimum_ondemand", &jsonschema.Schema{
		Type: "integer",
	})

	properties.Set("minimum_ondemand_percent", &jsonschema.Schema{
		Type: "integer",
	})

	return &jsonschema.Schema{
		Description: "Spot Capacity Overrides",
		OneOf: []*jsonschema.Schema{
			{
				Type:       "object",
				Properties: properties,
			},
			{
				Type: "boolean",
			},
		},
	}
}

func (obj *SpotOverrides) MarshalJSON() ([]byte, error) {
	if !obj.Enabled {
		return []byte("false"), nil
	}

	type t SpotOverrides
	res, err := util.Jsonify(t(*obj))
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}
