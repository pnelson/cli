package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// A Command is a single command in a command line application.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	// The return value is the exit code passed to os.Exit.
	Run func(cmd *Command, args []string) int

	// Usage is the one-line usage message.
	// The first word in the line is taken to be the command name.
	Usage string

	// Short is the short description shown in the help output.
	Short string

	// Long is the long message shown in the `help <cmd>` output.
	Long string

	// Flags is a set of flags specific to this command.
	Flags flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	i := strings.Index(c.Usage, " ")
	if i >= 0 {
		return c.Usage[:i]
	}

	return c.Usage
}

// usage prints the command usage and exits with an error code.
func (c *Command) usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.Usage)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}
