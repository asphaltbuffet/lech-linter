package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/lech-linter/cmd"
)

func Test_NewRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	require.NotNil(t, rootCmd)

	rootCmd2 := cmd.NewManCmd()
	require.NotNil(t, rootCmd2)

	assert.NotEqual(t, rootCmd, rootCmd2)
}

func TestRootCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string // filename -> content
		args         []string          // filenames to pass as args (resolved to temp paths)
		errAssertion require.ErrorAssertionFunc
		wantOut      []string // substrings expected in output
	}{
		{
			name: "no misspellings",
			files: map[string]string{
				"clean.txt": "Hello Lechlitner, welcome.\n",
			},
			args:         []string{"clean.txt"},
			errAssertion: require.NoError,
			wantOut:      []string{},
		},
		{
			name: "file with misspelling",
			files: map[string]string{
				"bad.txt": "Hello Lechniter!\n",
			},
			args:         []string{"bad.txt"},
			errAssertion: require.Error,
			wantOut:      []string{`possible misspelling "Lechniter"`},
		},
		{
			name:         "nonexistent file",
			files:        map[string]string{},
			args:         []string{"does-not-exist.txt"},
			errAssertion: require.Error,
			wantOut:      []string{},
		},
		{
			name: "multiple files aggregated",
			files: map[string]string{
				"a.txt": "Lechniter here\n",
				"b.txt": "Lechlitnre there\n",
			},
			args:         []string{"a.txt", "b.txt"},
			errAssertion: require.Error,
			wantOut: []string{
				`possible misspelling "Lechniter"`,
				`possible misspelling "Lechlitnre"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600)
				require.NoError(t, err)
			}

			// Resolve arg filenames to full temp paths.
			args := make([]string, len(tt.args))
			for i, a := range tt.args {
				args[i] = filepath.Join(dir, a)
			}

			var out bytes.Buffer

			root := cmd.NewRootCmd()
			root.SetOut(&out)
			root.SetArgs(args)

			tt.errAssertion(t, root.Execute())

			if len(tt.wantOut) == 0 {
				require.Empty(t, out.String())
			}

			for _, s := range tt.wantOut {
				assert.Contains(t, out.String(), s)
			}
		})
	}
}
