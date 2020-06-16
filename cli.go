// Package cli provides structure for command line applications.
package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

var (
	sep    = "_"
	mapper = strings.NewReplacer(".", sep, "/", sep, "-", sep, ",", sep)
)

// CLI represents a command line application.
type CLI struct {
	name           string
	usage          Renderer
	scope          string
	flags          []*Flag
	flagsMap       map[string]*Flag
	commands       map[string]*Command
	version        string
	stdout         io.Writer
	stderr         io.Writer
	afterParse     Handler
	helpHandler    Handler
	defaultHandler Handler
}

// New returns a new CLI application.
func New(name string, usage Renderer, flags []*Flag, opts ...Option) *CLI {
	c := &CLI{
		name:     name,
		usage:    usage,
		flags:    flags,
		flagsMap: make(map[string]*Flag),
		commands: make(map[string]*Command),
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}
	for _, option := range opts {
		option(c)
	}
	if c.scope != "" && !strings.HasSuffix(c.scope, "/") {
		c.scope += "/"
	}
	if c.helpHandler == nil {
		c.helpHandler = c.defaultHelpHandler
	}
	if c.defaultHandler == nil {
		c.defaultHandler = c.defaultDefaultHandler
	}
	c.Add("help", c.helpHandler, nil)
	if c.version != "" {
		c.Add("version", c.versionHandler, nil)
	}
	return c
}

// Add adds a new command.
func (c *CLI) Add(name string, handler Handler, flags []*Flag, opts ...CommandOption) *Command {
	name = strings.ToLower(name)
	if handler == nil {
		panic(fmt.Errorf("cli: command '%s' has nil handler", name))
	}
	_, ok := c.commands[name]
	if ok {
		panic(fmt.Errorf("cli: duplicate command '%s'", name))
	}
	cmd := NewCommand(name, handler, flags, opts...)
	c.commands[name] = cmd
	if cmd.alias != "" {
		dup, ok := c.commands[cmd.alias]
		if ok {
			panic(fmt.Errorf("cli: duplicate command alias '%s' for '%s'", cmd.alias, dup.name))
		}
		c.commands[cmd.alias] = cmd
	}
	return cmd
}

// Run parses the command line arguments, starting with the
// program name, and dispatches to the appropriate handler.
func (c *CLI) Run(args []string) error {
	if args == nil {
		args = os.Args
	} else if len(args) == 0 {
		args = []string{c.name}
	}
	for len(args) > 0 && args[len(args)-1] == "" {
		args = args[:len(args)-1]
	}
	return c.run(args)
}

// run parses the root command and dispatches to the given subcommand.
func (c *CLI) run(args []string) error {
	args, err := c.parse(args, c.flags)
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return c.defaultHandler(args)
	}
	name := args[0]
	cmd, ok := c.commands[name]
	if !ok {
		return c.commandNotFound(name)
	}
	args, err = c.parse(args, cmd.flags)
	if err != nil {
		return err
	}
	if c.afterParse != nil && name != "help" && name != "version" {
		err = c.afterParse(args)
		if err != nil {
			return err
		}
	}
	return cmd.handler(args)
}

// parse processes args as flags until there are no longer flags.
func (c *CLI) parse(args []string, flags []*Flag) ([]string, error) {
	err := c.initFlags(flags)
	if err != nil {
		return nil, err
	}
	return Parse(args[1:], flags)
}

// initFlags populates the application flag map and
// initial values from environment variables.
func (c *CLI) initFlags(flags []*Flag) error {
	for _, f := range flags {
		_, ok := c.flagsMap[f.name]
		if ok {
			return fmt.Errorf("Duplicate flag '%s'.", f.name)
		}
		c.flagsMap[f.name] = f
		if f.alias != "" {
			_, ok := c.flagsMap[f.alias]
			if ok {
				return fmt.Errorf("Duplicate short flag '%s' for '%s'.", f.alias, f.name)
			}
			c.flagsMap[f.alias] = f
		}
		if f.envKey == "" {
			key := strings.ToUpper(c.name + "_" + f.name)
			f.envKey = mapper.Replace(key)
		}
	}
	return nil
}

// commandNotFound prints helpful usage information and suggestions.
func (c *CLI) commandNotFound(name string) error {
	c.Errorf("Unknown command '%s'.\n", name)
	c.Errorf("Run '%s help' for usage information.\n", c.name)
	similar := make([]string, 0)
	for _, cmd := range c.commands {
		distance := 0
		if !strings.HasPrefix(cmd.name, name) {
			distance = levenshtein(name, cmd.name)
		}
		if distance < similarThreshold {
			similar = append(similar, cmd.name)
		}
	}
	if len(similar) > 0 {
		sort.Strings(similar)
		c.Errorf("\nDid you mean?\n\n")
		for _, name := range similar {
			c.Errorf("    %s\n", name)
		}
		c.Errorf("\n")
	}
	return ErrExitFailure
}

// Printf writes to the configured stdout writer.
func (c *CLI) Printf(format string, args ...interface{}) {
	fmt.Fprintf(c.stdout, format, args...)
}

// Errorf writes to the configured stderr writer.
func (c *CLI) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(c.stderr, format, args...)
}

// defaultHelpHandler is the default handler for the help command.
func (c *CLI) defaultHelpHandler(args []string) error {
	if len(args) == 0 {
		return c.Usage(c.stdout, c.scope)
	}
	if len(args) != 1 {
		c.Errorf("Too many arguments given.\n")
		c.Errorf("Run '%s help' for usage information.\n", c.name)
		c.Errorf("Run '%s help [command]' for more information about a command.\n", c.name)
		return ErrExitFailure
	}
	name := args[0]
	return c.Usage(c.stdout, name)
}

// defaultDefaultHandler is the default handler for naked commands.
func (c *CLI) defaultDefaultHandler(args []string) error {
	err := c.Usage(c.stderr, "")
	if err != nil {
		return err
	}
	return ErrExitFailure
}

// versionHandler is the handler for the version command.
func (c *CLI) versionHandler(args []string) error {
	c.Printf("%s\n", c.version)
	return nil
}

// Parse parses flag definitions from the argument list. Flag parsing stops
// at the first non-flag argument, including single or double hyphens followed
// by whitespace or end of input.
func Parse(args []string, flags []*Flag) ([]string, error) {
	m := make(map[string]*Flag)
	for _, f := range flags {
		m[f.name] = f
		if f.alias != "" {
			m[f.alias] = f
		}
		value, ok := os.LookupEnv(f.envKey)
		if ok {
			f.Set(value)
		}
	}
	key := ""
	for arg := ""; len(args) > 0; {
		arg, args = args[0], args[1:]
		if arg == "-" || arg == "--" {
			args = append([]string{arg}, args...)
			break
		}
		if key != "" {
			f, ok := m[key]
			if !ok {
				return nil, ErrUndefinedFlag(key)
			}
			if !f.kind.HasArg() {
				key = ""
				args = append([]string{arg}, args...)
				f.Set("true")
				continue
			}
			if arg[0] == '-' {
				return nil, ErrRequiresArg(key)
			}
			key = ""
			f.Set(arg)
			continue
		}
		if arg == "" || arg[0] != '-' {
			args = append([]string{arg}, args...)
			break
		}
		if arg[1] == '-' {
			arg = arg[2:]
		} else {
			arg = arg[1:]
		}
		if !unicode.IsLetter(rune(arg[0])) {
			return nil, ErrFlagSyntax(arg)
		}
		i := strings.Index(arg, "=")
		if i == -1 {
			key = arg
		} else {
			key = arg[:i]
			f, ok := m[key]
			if !ok {
				return nil, ErrUndefinedFlag(key)
			}
			key = ""
			f.Set(arg[i+1:])
		}
	}
	if key != "" {
		f, ok := m[key]
		if ok {
			if f.kind.HasArg() {
				return nil, ErrRequiresArg(key)
			}
			f.Set("true")
		} else {
			return nil, ErrUndefinedFlag(key)
		}
	}
	return args, nil
}
