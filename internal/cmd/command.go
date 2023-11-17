package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

type (
	BuildCommand struct {
		cmd  *cobra.Command
		opts *BuildOptions
	}

	BuildOptions struct {
		config        string
		ids           []string
		logDirFlag    bool
		configStrFlag bool
		configFlag    bool
		scriptFlag    bool
		initFlag      bool
		timeout       time.Duration
	}
)

func NewBuildCommand(name string) *BuildCommand {
	return &BuildCommand{
		cmd: &cobra.Command{
			Use: name,
		},
	}
}

func (c *BuildCommand) AddCommandMessage(shortMessage, longMessage string) *BuildCommand {
	c.cmd.Short = shortMessage
	c.cmd.Long = longMessage
	return c
}

func (c *BuildCommand) AddCommandRun(run func(cmd *cobra.Command, args []string)) *BuildCommand {
	c.cmd.Run = run
	return c
}

func (c *BuildCommand) AddCommandFlag(name, defaultValue, description string) *BuildCommand {
	c.cmd.Flags().StringVar(&c.opts.config, name, defaultValue, description)
	return c
}

func (c *BuildCommand) AddCommandFlagArray(name string, defaultValue []string, description string) *BuildCommand {
	c.cmd.Flags().StringSliceVar(&c.opts.ids, name, defaultValue, description)
	return c
}

func (c *BuildCommand) AddCommandFlagBool(name string, defaultValue bool, description string) *BuildCommand {
	c.cmd.Flags().BoolVar(&c.opts.logDirFlag, name, defaultValue, description)
	return c
}

func (c *BuildCommand) AddCommandFlagBoolVar(name string, defaultValue bool, description string) *BuildCommand {
	c.cmd.Flags().BoolVar(&c.opts.configStrFlag, name, defaultValue, description)
	return c
}

func (c *BuildCommand) AddCommandFlagBoolVarP(name string, shortName string, defaultValue bool, description string) *BuildCommand {
	c.cmd.Flags().BoolVarP(&c.opts.configFlag, name, shortName, defaultValue, description)
	return c
}

func (c *BuildCommand) SilentUsage() *BuildCommand {
	c.cmd.SilenceUsage = true
	return c
}

func (c *BuildCommand) SilentErrors() *BuildCommand {
	c.cmd.SilenceErrors = true
	return c
}

func (c *BuildCommand) InitDefaultHelpFlag() *BuildCommand {
	c.cmd.InitDefaultHelpFlag()
	return c
}

func (c *BuildCommand) Validate() error {
	return c.cmd.ValidateArgs(c.opts.ids)
}

func (c *BuildCommand) BuildCommand() *cobra.Command {
	return c.cmd
}
