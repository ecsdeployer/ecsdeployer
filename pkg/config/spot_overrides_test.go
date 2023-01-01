package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSpotOverrides(t *testing.T) {

	// For Spot/Ondemand the array is [base, weight]

	tables := []struct {
		disabled     bool
		expectedSpot []int32
		expectedOD   []int32
		obj          *config.SpotOverrides
	}{
		// {true, nil, nil, config.SpotOverrides{}},
		{true, nil, []int32{0, 1}, &config.SpotOverrides{}},
		{true, nil, []int32{0, 1}, config.NewSpotOnDemand()},
		{false, []int32{0, 100}, []int32{0, 1}, &config.SpotOverrides{Enabled: true, MinimumOnDemand: util.Ptr[int32](0), MinimumOnDemandPercent: util.Ptr[int32](1)}},
		{false, []int32{0, 1}, nil, &config.SpotOverrides{Enabled: true}},
		{false, []int32{0, 100}, []int32{1, 0}, &config.SpotOverrides{Enabled: true, MinimumOnDemand: util.Ptr[int32](1), MinimumOnDemandPercent: util.Ptr[int32](0)}},
	}

	for entryNum, table := range tables {
		t.Run(fmt.Sprintf("entry_%02d", entryNum), func(t *testing.T) {
			obj := table.obj
			obj.ApplyDefaults()

			if table.expectedOD == nil && table.expectedSpot == nil {
				t.Fatal("TEST BROKEN! you cant have SPOT and OD set to nil, that doesnt make sense")
			}

			if err := obj.Validate(); err != nil {
				require.NoErrorf(t, err, "Entry<%d> DID NOT PASS VALIDATION! TEST IS BROKE err: %v", entryNum, err)
			}

			require.Equalf(t, table.disabled, obj.IsDisabled(), "IsDisabled")

			expectedEntries := 0
			if table.expectedOD != nil {
				expectedEntries += 1

				require.Truef(t, obj.WantsOnDemand(), "WantsOnDemand")

			}
			if table.expectedSpot != nil {
				expectedEntries += 1

				// they expected to not have ondemand, but it says they do
				if expectedEntries == 1 && obj.WantsOnDemand() {
					require.Falsef(t, obj.WantsOnDemand(), "Expected WantsOnDemand()==false but it wasnt")
				}
			}

			strategy := obj.ExportCapacityStrategy()
			ebStrategy := obj.ExportCapacityStrategyEventBridge()

			require.Equalf(t, len(strategy), len(ebStrategy), "export strategy mismatch")

			require.Lenf(t, strategy, expectedEntries, "num exports mismatch")

			// ensure they are both identical, and only differ because AWS loves types
			for i := range strategy {

				require.Equalf(t, *strategy[i].CapacityProvider, *ebStrategy[i].CapacityProvider, "CapacityProvider")
				require.Equalf(t, strategy[i].Base, ebStrategy[i].Base, "Base")
				require.Equalf(t, strategy[i].Weight, ebStrategy[i].Weight, "Weight")
			}

			stratMap := make(map[string][]int32, len(strategy))

			for _, entry := range strategy {
				stratMap[*entry.CapacityProvider] = []int32{
					entry.Base,
					entry.Weight,
				}
			}

			if table.expectedOD != nil {
				ent := stratMap["FARGATE"]

				require.Equalf(t, table.expectedOD[0], ent[0], "FARGATE_base entry=%d", entryNum)
				require.Equalf(t, table.expectedOD[1], ent[1], "FARGATE_weight entry=%d", entryNum)
			}

			if table.expectedSpot != nil {
				ent := stratMap["FARGATE_SPOT"]

				require.Equalf(t, table.expectedSpot[0], ent[0], "FARGATE_SPOT_base entry=%d", entryNum)
				require.Equalf(t, table.expectedSpot[1], ent[1], "FARGATE_SPOT_weight entry=%d", entryNum)
			}
		})
	}
}

func TestSpotOverrides_Marshalling(t *testing.T) {
	tables := []struct {
		str        string
		failure    bool
		exp        *config.SpotOverrides
		expMarshal string
	}{

		{
			str:        `false`,
			exp:        &config.SpotOverrides{Enabled: false},
			expMarshal: "false",
		},

		{
			str:        `true`,
			exp:        &config.SpotOverrides{Enabled: true},
			expMarshal: `{"enabled":true}`,
		},
		{
			str:        "enabled: true\nminimum_ondemand: 1",
			exp:        &config.SpotOverrides{Enabled: true, MinimumOnDemand: util.Ptr[int32](1)},
			expMarshal: `{"enabled":true,"minimum_ondemand":1}`,
		},
	}

	for i, table := range tables {
		t.Run(fmt.Sprintf("entry_%02d", i), func(t *testing.T) {
			obj, err := yaml.ParseYAMLString[config.SpotOverrides](table.str)
			if table.failure {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, obj)

			require.EqualValuesf(t, table.exp.Enabled, obj.Enabled, "Enabled")
			require.EqualValuesf(t, table.exp.MinimumOnDemand, obj.MinimumOnDemand, "MinimumOnDemand")
			require.EqualValuesf(t, table.exp.MinimumOnDemandPercent, obj.MinimumOnDemandPercent, "MinimumOnDemandPercent")

			if table.expMarshal != "" {
				jsonData, err := json.Marshal(obj)
				require.NoError(t, err)

				require.Equal(t, table.expMarshal, string(jsonData))
			}

		})
	}
}
