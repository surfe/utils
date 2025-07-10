package utils

func AppendWithoutDuplicate(a []string, v string) []string {
	for _, x := range a {
		if x == v {
			return a
		}
	}

	return append(a, v)
}

func MergeWithoutDuplicates(a []string, b []string) []string {
	for _, x := range b {
		a = AppendWithoutDuplicate(a, x)
	}

	return a
}

func ChunkBy[T any](items []T, chunkSize int) [][]T {
	var chunks [][]T
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

func SliceFromIndexOrNil[T any](a []T, i int) []T {
	if i >= 0 && i < len(a) {
		return a[i:]
	}

	return nil
}

func SafeAppendToSlice[T any](slicePtr *[]T, item T) *[]T {
	if slicePtr == nil {
		slicePtr = &[]T{}
	}

	*slicePtr = append(*slicePtr, item)

	return slicePtr
}
