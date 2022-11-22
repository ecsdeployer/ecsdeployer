package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSpotOverrides(t *testing.T) {

	// For Spot/Ondemand the array is [base, weight]

	tables := []struct {
		disabled     bool
		expectedSpot []int32
		expectedOD   []int32
		obj          config.SpotOverrides
	}{
		// {true, nil, nil, config.SpotOverrides{}},
		{true, nil, []int32{0, 1}, config.SpotOverrides{}},
		{false, []int32{0, 100}, []int32{0, 1}, config.SpotOverrides{Enabled: true, MinimumOnDemand: util.Ptr[int32](0), MinimumOnDemandPercent: util.Ptr[int32](1)}},
		{false, []int32{0, 1}, nil, config.SpotOverrides{Enabled: true}},
		{false, []int32{0, 100}, []int32{1, 0}, config.SpotOverrides{Enabled: true, MinimumOnDemand: util.Ptr[int32](1), MinimumOnDemandPercent: util.Ptr[int32](0)}},
	}

	for entryNum, table := range tables {
		obj := table.obj
		obj.ApplyDefaults()

		if table.expectedOD == nil && table.expectedSpot == nil {
			t.Fatal("TEST BROKEN! you cant have SPOT and OD set to nil, that doesnt make sense")
		}

		if err := obj.Validate(); err != nil {
			require.NoErrorf(t, err, "Entry<%d> DID NOT PASS VALIDATION! TEST IS BROKE err: %v", entryNum, err)
		}

		if obj.IsDisabled() != table.disabled {
			t.Errorf("Entry<%d> Expected IsDisabled()==%t but it wasnt", entryNum, table.disabled)
		}

		expectedEntries := 0
		if table.expectedOD != nil {
			expectedEntries += 1

			if !obj.WantsOnDemand() {
				t.Errorf("Entry<%d> Expected WantsOnDemand()==true but it wasnt", entryNum)
			}

		}
		if table.expectedSpot != nil {
			expectedEntries += 1

			// they expected to not have ondemand, but it says they do
			if expectedEntries == 1 && obj.WantsOnDemand() {
				t.Errorf("Entry<%d> Expected WantsOnDemand()==false but it wasnt", entryNum)
			}
		}

		strategy := obj.ExportCapacityStrategy()
		ebStrategy := obj.ExportCapacityStrategyEventBridge()

		if len(strategy) != len(ebStrategy) {
			t.Errorf("Entry<%d> export strategy mismatch", entryNum)
		}

		if len(strategy) != expectedEntries {
			t.Errorf("Entry<%d> num exports mismatch. expected=%d got=%d", entryNum, expectedEntries, len(strategy))
		}

		// ensure they are both identical, and only differ because AWS loves types
		for i := range strategy {
			if *strategy[i].CapacityProvider != *ebStrategy[i].CapacityProvider {
				t.Errorf("Entry<%d>::<%d> Expected CapacityProvider==%v but got %v", entryNum, i, *strategy[i].CapacityProvider, *ebStrategy[i].CapacityProvider)
				continue
			}

			if strategy[i].Base != ebStrategy[i].Base {
				t.Errorf("Entry<%d>::<%d> Expected Base==%v but got %v", entryNum, i, strategy[i].Base, ebStrategy[i].Base)
				continue
			}

			if strategy[i].Weight != ebStrategy[i].Weight {
				t.Errorf("Entry<%d>::<%d> Expected Weight==%v but got %v", entryNum, i, strategy[i].Weight, ebStrategy[i].Weight)
				continue
			}
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
			if table.expectedOD[0] != ent[0] {
				t.Errorf("Entry<%d> expected to have FARGATE base=%d but got %d", entryNum, table.expectedOD[0], ent[0])
			}
			if table.expectedOD[1] != ent[1] {
				t.Errorf("Entry<%d> expected to have FARGATE weight=%d but got %d", entryNum, table.expectedOD[1], ent[1])
			}
		}

		if table.expectedSpot != nil {
			ent := stratMap["FARGATE_SPOT"]
			if table.expectedSpot[0] != ent[0] {
				t.Errorf("Entry<%d> expected to have FARGATE_SPOT base=%d but got %d", entryNum, table.expectedSpot[0], ent[0])
			}
			if table.expectedSpot[1] != ent[1] {
				t.Errorf("Entry<%d> expected to have FARGATE_SPOT weight=%d but got %d", entryNum, table.expectedSpot[1], ent[1])
			}
		}
	}
}
