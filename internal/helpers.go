package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FindBestMatch[T any](items map[string]T, getTitle func(T) string, query string) (T, bool) {
	var bestItem T
	highScore := -1
	for _, item := range items {
		s := Score(getTitle(item), query)
		if s > highScore && s > 0 {
			highScore = s
			bestItem = item
		}
	}
	if highScore == -1 {
		return bestItem, false
	}
	return bestItem, true
}

func Confirm(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}
