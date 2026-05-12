package internal

import (
	"strings"
	"unicode"
)

func normalize(s string) string {
	s = strings.ToLower(s)

	var b strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func Score(title, query string) int {
	title = normalize(title)
	query = normalize(query)

	// Quick exit for exact matches
	if title == query {
		return 1000
	}

	// Quick exit for substring containment
	if strings.Contains(title, query) {
		return 500 + len(query)
	}

	titleWords := strings.Fields(title)
	queryWords := strings.Fields(query)

	// Use a map for O(1) lookups instead of nested loops
	titleMap := make(map[string]int)
	for i, word := range titleWords {
		titleMap[word] = i
	}

	currentScore := 0
	for i, qw := range queryWords {
		if pos, exists := titleMap[qw]; exists {
			currentScore += 10

			// Proximity bonus: if the next word in query also matches
			// the next word in the title
			if i+1 < len(queryWords) && pos+1 < len(titleWords) {
				if queryWords[i+1] == titleWords[pos+1] {
					currentScore += 50
				}
			}
		}
	}

	return currentScore
}
