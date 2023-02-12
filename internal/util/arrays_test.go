package util

import (
	"fmt"
	"testing"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

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
		t.Run(fmt.Sprintf("test_%02d", i+1), func(t *testing.T) {
			chunks := ChunkArray(table.input, table.chunkSize)

			require.Lenf(t, chunks, table.expChunks, "expected %d chunks, but got %d", table.expChunks, len(chunks))

			totalItems := 0

			lastVal, chunks := chunks[len(chunks)-1], chunks[:len(chunks)-1]

			for c, chunk := range chunks {
				chunkLen := len(chunk)
				totalItems += chunkLen
				require.Equalf(t, table.chunkSize, chunkLen, "expected chunk <%d> to have %d items but it only had %d", c, table.chunkSize, chunkLen)
			}

			require.LessOrEqual(t, len(lastVal), table.chunkSize, "last chunk has too many items")
			totalItems += len(lastVal)

			require.Equal(t, len(table.input), totalItems, "missing items??")
		})

	}
}

func TestStrArrayToInterArray(t *testing.T) {

	input := ecsTypes.TaskDefinitionStatusActive.Values()
	inters := StrArrayToInterArray(input)

	for i, val := range input {
		require.EqualValues(t, val, inters[i])
	}

	// require.IsType(t, []any, inters)

}
