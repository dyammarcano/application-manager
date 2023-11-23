package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type (
	BuildCommand struct {
		Cmd     *cobra.Command
		Options []State
	}

	State struct {
		Touched any
		Flag    Flag
	}

	Flag struct {
		Name, Description string
		Default           any
		Value             any
	}
)

func NewCommandBuilder(name string) *BuildCommand {
	return &BuildCommand{
		Options: make([]State, 0),
		Cmd: &cobra.Command{
			Use: name,
		},
	}
}

func (c *BuildCommand) AddCommand(buildCommand *BuildCommand) *BuildCommand {
	c.Cmd.AddCommand(buildCommand.Cmd)
	// add options to parent command
	for k, v := range buildCommand.Options {
		c.Options[k] = v
	}
	return c
}

func (c *BuildCommand) AddCommandShortMessage(shortMessage string) *BuildCommand {
	c.Cmd.Short = shortMessage
	return c
}

func (c *BuildCommand) AddCommandLongMessage(longMessage string) *BuildCommand {
	c.Cmd.Long = longMessage
	return c
}

func (c *BuildCommand) AddCommandRun(run func(cmd *cobra.Command, args []string)) *BuildCommand {
	c.Cmd.Run = run
	return c
}

func (c *BuildCommand) AddCommandFlag(name string, defaultValue any, description string) *BuildCommand {
	c.addCommandFlag(name, defaultValue, description, false)
	if err := c.Cmd.Flags().MarkHidden(name); err != nil {
		c.Cmd.Printf("Error binding flag: %s\n", err)
		os.Exit(1)
	}
	return c
}

func (c *BuildCommand) addCommandFlag(name string, defaultValue any, description string, persistent bool) *BuildCommand {
	state := State{
		Touched: true,
		Flag: Flag{
			Name:        name,
			Description: description,
			Default:     defaultValue,
			Value:       defaultValue,
		},
	}

	c.Options = append(c.Options, state)

	switch v := defaultValue.(type) {
	case bool:
		if persistent {
			c.Cmd.PersistentFlags().Bool(name, v, description)
		} else {
			c.Cmd.Flags().Bool(name, v, description)
		}
	case string:
		if persistent {
			c.Cmd.PersistentFlags().String(name, v, description)
		} else {
			c.Cmd.Flags().String(name, v, description)
		}
	case int, int8, int16, int32, int64:
		if persistent {
			c.Cmd.PersistentFlags().Int64(name, v.(int64), description)
		} else {
			c.Cmd.Flags().Int64(name, v.(int64), description)
		}
	default:
		fmt.Printf("Invalid type: %s\n", v)
		os.Exit(1)
	}

	return c
}

func (c *BuildCommand) AddCommandFlagPersistent(name string, defaultValue any, description string) *BuildCommand {
	c.addCommandFlag(name, defaultValue, description, true)
	if err := c.Cmd.PersistentFlags().MarkHidden(name); err != nil {
		c.Cmd.Printf("Error binding flag: %s\n", err)
		os.Exit(1)
	}

	return c
}

func (c *BuildCommand) SilentUsage() *BuildCommand {
	c.Cmd.SilenceUsage = true
	return c
}

func (c *BuildCommand) SilentErrors() *BuildCommand {
	c.Cmd.SilenceErrors = true
	return c
}

func (c *BuildCommand) InitDefaultHelpFlag() *BuildCommand {
	c.Cmd.InitDefaultHelpFlag()
	return c
}

func (c *BuildCommand) Build() *BuildCommand {
	return c
}
