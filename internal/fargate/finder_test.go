package fargate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindFargateBestFit(t *testing.T) {

	smallest := DefaultFargateResources[0]
	biggest := DefaultFargateResources[len(DefaultFargateResources)-1]

	tables := []struct {
		cpu    int32
		mem    int32
		expCpu int32
		expMem int32
	}{
		{0, 0, smallest.Cpu, smallest.Memory},
		{1000, 3000, 1024, 3072},
		{99999999, 0, biggest.Cpu, biggest.Memory},
		{0, 99999999, biggest.Cpu, biggest.Memory},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("FF_cpu_%d__memory_%d", table.cpu, table.mem), func(t *testing.T) {

			res := FindFargateBestFit(table.cpu, table.mem)
			require.Equalf(t, table.expCpu, res.Cpu, "CPU")
			require.Equalf(t, table.expMem, res.Memory, "Memory")
		})

	}
}

func TestExceedsLargest(t *testing.T) {
	tables := []struct {
		cpu int32
		mem int32
		ret bool
	}{
		{128, 128, false},
		{9999999, 128, true},
		{128, 99999999, true},
	}

	for _, table := range tables {
		ret := ExceedsLargest(table.cpu, table.mem)
		require.Equal(t, table.ret, ret)
	}
}

func TestFindFargateBestFitOrTrust(t *testing.T) {
	tables := []struct {
		cpu         int32
		mem         int32
		expectedCpu int32
		expectedMem int32
	}{
		{128, 128, 256, 512},
		{16384, 122880, 16384, 122880},
		{16384, 120000, 16384, 122880},
		{32768, 122880, 32768, 122880},
	}

	for _, table := range tables {
		ret := FindFargateBestFitOrTrust(table.cpu, table.mem)
		require.Equal(t, table.expectedCpu, ret.Cpu)
		require.Equal(t, table.expectedMem, ret.Memory)
	}
}
