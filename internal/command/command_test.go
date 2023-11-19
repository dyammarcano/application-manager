package command

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBuildCommand(t *testing.T) {
	cmd := NewCommandBuilder("test").
		AddCommandShortMessage("short test").
		AddCommandLongMessage("long test").
		AddCommandRun(func(cmd *cobra.Command, args []string) {}).
		AddCommandFlag("config", "", "config file").
		Build()
	assert.Nilf(t, cmd, "BuildOptions should be nil")
}
