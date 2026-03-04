package checker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/lech-linter/internal/checker"
)

func TestLevenshtein(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"identical - Proper", "Lechlitner", "Lechlitner", 0},
		{"identical - lower", "lechlitner", "lechlitner", 0},
		{"wrong letter", "lechlitner", "xechlitner", 1},
		{"2 edits", "lechlitner", "lechniter", 2},
		{"both empty", "", "", 0},
		{"empty check", "abc", "", 3},
		{"empty src", "", "abc", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := checker.Levenshtein(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsMisspelling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		word          string
		boolAssertion assert.BoolAssertionFunc
	}{
		{"proper", "Lechlitner", assert.False},
		{"all lower", "lechlitner", assert.False},
		{"all caps", "LECHLITNER", assert.False},
		{"missing + swap", "Lechniter", assert.True},
		{"end swapped", "Lechlitnre", assert.True},
		{"too short", "hello", assert.False},
		{"plural", "Lechlitners", assert.True},
		{"very plural", "Lechlitnerss", assert.True},
		{"too long", "superlongword", assert.False},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := checker.IsMisspelling(tt.word)
			tt.boolAssertion(t, got)
		})
	}
}
