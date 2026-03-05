package checker

import "strings"

const (
	target    = "lechlitner"
	threshold = 3
	minLen    = len(target) - threshold
	maxLen    = len(target) + threshold
)

// Levenshtein returns the edit distance between two strings.
func Levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	aLen, bLen := len(ra), len(rb)

	if aLen == 0 {
		return bLen
	}

	if bLen == 0 {
		return aLen
	}

	prev := make([]int, bLen+1)
	cur := make([]int, bLen+1)

	for i := range prev {
		prev[i] = i
	}

	for i := 1; i <= aLen; i++ {
		cur[0] = i
		for j := 1; j <= bLen; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(
				cur[j-1]+1,
				prev[j]+1,
				prev[j-1]+cost,
			)
		}
		prev, cur = cur, prev
	}

	return prev[bLen]
}

// IsMisspelling returns true if word is within Levenshtein distance of
// "Lechlitner" but is not the correct spelling.
func IsMisspelling(word string) bool {
	lower := strings.ToLower(word)

	// Length filter: skip words clearly too short or too long.
	// use of []rune() here is intentional to avoid issues with bytes vs runes
	if len([]rune(lower)) < minLen || len([]rune(lower)) > maxLen {
		return false
	}

	if lower == target {
		return false
	}

	return Levenshtein(lower, target) <= threshold
}
