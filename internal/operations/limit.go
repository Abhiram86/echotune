package operations

func Limit[T any](list []T, limit int) []T {
	if limit <= 0 || limit >= len(list) {
		return list
	}
	return list[:limit]
}
