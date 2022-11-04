package util

func ChunkArray[T any](input []T, chunkSize int) [][]T {
	// if len(input) <= chunkSize {
	// 	return [][]T{input}
	// }
	// batches := make([][]T, 0, (len(input)+chunkSize-1)/chunkSize)
	// for chunkSize < len(input) {
	// 	input, batches = input[chunkSize:], append(batches, input[0:chunkSize:chunkSize])
	// }
	// return batches
	if len(input) <= chunkSize {
		return [][]T{input}
	}

	chunks := make([][]T, 0, (len(input)+chunkSize-1)/chunkSize)
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond cap
		if end > len(input) {
			end = len(input)
		}

		chunks = append(chunks, input[i:end])
	}
	return chunks
}

// This is used to allow the AWS Enum string lists to be cast to an array of interfaces
// used for JSONSchema
func StrArrayToInterArray[T ~string](things []T) []interface{} {
	arr := make([]interface{}, len(things))
	for i, val := range things {
		arr[i] = string(val)
	}
	return arr
}
