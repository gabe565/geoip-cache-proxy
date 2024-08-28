package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithVersion(t *testing.T) {
	cmd := &cobra.Command{}
	const expect = "1.0.0"
	WithVersion(expect)(cmd)
	require.NotNil(t, cmd.Annotations)
	assert.NotEmpty(t, cmd.Version)
	assert.Equal(t, expect, cmd.Version)
	assert.NotEmpty(t, cmd.Annotations[VersionKey])
}
