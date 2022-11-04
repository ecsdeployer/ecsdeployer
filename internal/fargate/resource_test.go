package fargate

import "testing"

func TestFargateResourceStrings(t *testing.T) {
	tables := []struct {
		res FargateResource
		cpu string
		mem string
	}{
		{FargateResource{1024, 2048}, "1024", "2048"},
	}

	for _, table := range tables {
		if table.res.CpuString() != table.cpu {
			t.Errorf("expected cpu %s to be %s but it wasnt", table.res.CpuString(), table.cpu)
		}

		if table.res.MemoryString() != table.mem {
			t.Errorf("expected mem %s to be %s but it wasnt", table.res.MemoryString(), table.mem)
		}
	}
}
