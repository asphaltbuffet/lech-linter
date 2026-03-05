package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/lech-linter/internal/version"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	// just make sure there's something in version as a default
	assert.NotEmpty(t, version.Version)
}
