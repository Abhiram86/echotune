package operations

import "sort"

func Sort[T any](list []T, less func(a, b T) bool) []T {
	sort.Slice(list, func(i, j int) bool {
		return less(list[i], list[j])
	})
	return list
}
