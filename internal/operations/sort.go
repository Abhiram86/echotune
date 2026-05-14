package operations

import "sort"

func Sort[T any](list []T, less func(a, b T) bool) []T {
	sort.Slice(list, func(i, j int) bool {
		return less(list[i], list[j])
	})
	return list
}

func ToSortedSlice[T any](m map[string]T, less func(a, b *T) bool) []T {
	slice := make([]T, 0, len(m))
	for _, v := range m {
		slice = append(slice, v)
	}
	sort.Slice(slice, func(i, j int) bool {
		return less(&slice[i], &slice[j])
	})
	return slice
}
