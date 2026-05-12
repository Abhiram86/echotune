package operations

import "math/rand"

func Shuffle[T any](list []T) []T {
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
	return list
}
