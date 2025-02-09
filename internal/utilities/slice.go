package utilities

func Filter[T any, Slice ~[]T](
	collection Slice,
	predicate func(item T, index int) bool,
) Slice {
	result := make(Slice, 0, len(collection))

	for i := range collection {
		if predicate(collection[i], i) {
			result = append(result, collection[i])
		}
	}

	return result
}
