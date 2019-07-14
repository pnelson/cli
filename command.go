package cli

import "strings"

// Command represents an application command.
type Command struct {
	name    string
	alias   string
	usage   string
	flags   []*Flag
	handler Handler
}

// Handler represents a command handler.
type Handler func(args []string) error

// NewCommand returns a new command.
func NewCommand(name string, handler Handler, usage string, flags []*Flag, opts ...CommandOption) *Command {
	c := &Command{
		name:    name,
		usage:   strings.TrimSpace(usage),
		flags:   flags,
		handler: handler,
	}
	for _, option := range opts {
		option(c)
	}
	return c
}

// CommandOption represents a functional option for command configuration.
type CommandOption func(*Command)

// Alias sets the command alias.
func Alias(name string) CommandOption {
	return func(c *Command) {
		c.alias = name
	}
}
