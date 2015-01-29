/*
Package cli provides structure for command line applications with sub-commands.

This package was heavily influenced by the Go command line application.
Commands help and version are implemented by default. The usage information
is pretty printed in an opinionated format.
*/
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// An Application represents a command line application.
type Application struct {
	name     string
	version  string
	commands []*Command
}

const similarThreshold = 5

// New creates a basic Application with help and version commands.
func New(name, version string) *Application {
	app := &Application{
		name:    name,
		version: version,
	}

	app.Command(&Command{
		native: true,
		Usage:  "help",
		Short:  "Output this usage information.",
	})

	app.Command(&Command{
		native: true,
		Usage:  "version",
		Short:  "Output the application version.",
		Run: func(args []string) int {
			fmt.Printf("%s v%s\n", app.name, app.version)
			return 0
		},
	})

	return app
}

// Command registers a command with the application.
func (a *Application) Command(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// Run will parse flags and dispatch to the command.
func (a *Application) Run() {
	flag.Usage = a.usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		a.usage()
	}

	if args[0] == "help" {
		a.help(args[1:])
		return
	}

	for _, cmd := range a.commands {
		name := cmd.Name()
		if name != args[0] || cmd.Run == nil {
			if strings.HasPrefix(name, args[0]) {
				cmd.distance = 0
			} else {
				cmd.distance = levenshtein(args[0], name)
			}
			continue
		}

		cmd.Flags.Usage = cmd.usage
		cmd.Flags.Parse(args[1:])
		args = cmd.Flags.Args()
		code := cmd.Run(args)
		os.Exit(code)
	}

	sort.Sort(byDistance(a.commands))

	fmt.Fprintf(os.Stderr, "%s: unknown command %#q\n", a.name, args[0])
	fmt.Fprintf(os.Stderr, "Run `%s help` for usage.\n", a.name)

	similar := a.similar()
	if similar != nil {
		fmt.Fprintf(os.Stderr, "\nDid you mean?\n")
		for _, cmd := range similar {
			fmt.Fprintf(os.Stderr, "    %s\n", cmd.Name())
		}
	}

	os.Exit(2)
}

func (a *Application) help(args []string) {
	if len(args) == 0 {
		a.printUsage(os.Stdout)
		return
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s help [command]\n", a.name)
		fmt.Fprintf(os.Stderr, "Too many arguments given.\n")
		os.Exit(2)
	}

	name := args[0]
	for _, cmd := range a.commands {
		if cmd.Name() != name {
			continue
		}

		data := struct {
			Name    string
			Command *Command
		}{
			Name:    a.name,
			Command: cmd,
		}

		tmpl(os.Stdout, helpTemplate, &data)
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.\n", name)
	fmt.Fprintf(os.Stderr, "Run `%s help`.\n", a.name)
	os.Exit(2)
}

func (a *Application) similar() []*Command {
	var rv []*Command
	for _, cmd := range a.commands {
		if !cmd.native && cmd.distance < similarThreshold {
			rv = append(rv, cmd)
		}
	}
	return rv
}

func (a *Application) printUsage(w io.Writer) {
	data := struct {
		Name     string
		Commands []*Command
	}{
		Name:     a.name,
		Commands: a.commands,
	}

	tmpl(w, usageTemplate, &data)
}

func (a *Application) usage() {
	a.printUsage(os.Stderr)
	os.Exit(2)
}
