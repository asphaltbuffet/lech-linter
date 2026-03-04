package checker_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/lech-linter/internal/checker"
)

func TestScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		filename string
		want     []checker.Finding
	}{
		{
			name:     "no matches",
			input:    "Hello stephanie, welcome.\n",
			filename: "test.txt",
			want:     nil,
		},
		{
			name:     "no misspellings",
			input:    "Hello Lechlitner, welcome.\n",
			filename: "test.txt",
			want:     nil,
		},
		{
			name:     "one misspelling",
			input:    "Hello Lechniter!\n",
			filename: "test.txt",
			want: []checker.Finding{
				{File: "test.txt", Line: 1, Col: 7, Word: "Lechniter"},
			},
		},
		{
			name:     "mixed multiline",
			input:    "Hello Lechniter!\nNothing is spelled wrong here...\nLechlitner spelled right\nLechlinter is wrong\n",
			filename: "test.txt",
			want: []checker.Finding{
				{File: "test.txt", Line: 1, Col: 7, Word: "Lechniter"},
				{File: "test.txt", Line: 4, Col: 1, Word: "Lechlinter"},
			},
		},
		{
			name:     "misspelling with punctuation",
			input:    "Lechlitnre, please.\n",
			filename: "f.txt",
			want: []checker.Finding{
				{File: "f.txt", Line: 1, Col: 1, Word: "Lechlitnre"},
			},
		},
		{
			name:     "multiple lines",
			input:    "ok\nLechniter here\n",
			filename: "multi.txt",
			want: []checker.Finding{
				{File: "multi.txt", Line: 2, Col: 1, Word: "Lechniter"},
			},
		},
		{
			name:     "stdin label",
			input:    "Lechniter\n",
			filename: "<stdin>",
			want: []checker.Finding{
				{File: "<stdin>", Line: 1, Col: 1, Word: "Lechniter"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := strings.NewReader(tt.input)
			got, err := checker.Scan(r, tt.filename)
			require.NoError(t, err)

			require.Len(t, got, len(tt.want))
			assert.Equal(t, tt.want, got)
		})
	}
}
