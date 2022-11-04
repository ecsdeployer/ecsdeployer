package fargate

import "testing"

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
		res := FindFargateBestFit(table.cpu, table.mem)

		if res.Cpu != table.expCpu || res.Memory != table.expMem {
			t.Errorf("FF(%d, %d) expected (%d,%d) but got (%d,%d)",
				table.cpu, table.mem,
				table.expCpu, table.expMem,
				res.Cpu, res.Memory,
			)
		}

	}
}
