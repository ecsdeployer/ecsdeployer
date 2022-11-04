package util

import "testing"

func TestChunkArray(t *testing.T) {
	tables := []struct {
		input     []int
		chunkSize int
		expChunks int
	}{
		{make([]int, 30), 10, 3},
		{make([]int, 5), 10, 1},
		{make([]int, 9), 10, 1},
		{make([]int, 10), 10, 1},
		{make([]int, 31), 10, 4},
	}

	for i, table := range tables {
		chunks := ChunkArray(table.input, table.chunkSize)

		t.Logf("Entry<%d> == %v", i, chunks)

		if len(chunks) != table.expChunks {
			t.Fatalf("Entry<%d> expected %d chunks, but got %d", i, table.expChunks, len(chunks))
		}

		totalItems := 0

		lastVal, chunks := chunks[len(chunks)-1], chunks[:len(chunks)-1]

		for c, chunk := range chunks {
			chunkLen := len(chunk)
			totalItems += chunkLen
			if chunkLen != table.chunkSize {
				t.Fatalf("Entry<%d> expected chunk <%d> to have %d items but it only had %d", i, c, table.chunkSize, chunkLen)
			}
		}

		if len(lastVal) > table.chunkSize {
			t.Fatalf("Entry<%d> last chunk has too many items. (had %d)", i, len(lastVal))
		}

		totalItems += len(lastVal)

		if totalItems != len(table.input) {
			t.Fatalf("Entry<%d> Has missing items??. (start=%d, now=%d)", i, len(table.input), totalItems)
		}

	}
}
