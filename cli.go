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
)

// An Application represents a command line application.
type Application struct {
	name     string
	version  string
	commands []*Command
}

// New creates a basic Application with help and version commands.
func New(name, version string) *Application {
	app := &Application{
		name:    name,
		version: version,
	}

	app.Command(&Command{
		Usage: "help",
		Short: "Output this usage information.",
	})

	app.Command(&Command{
		Usage: "version",
		Short: "Output the application version.",
		Run: func(cmd *Command, args []string) int {
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
		if cmd.Name() != args[0] || cmd.Run == nil {
			continue
		}

		cmd.Flags.Usage = cmd.usage
		cmd.Flags.Parse(args[1:])
		args = cmd.Flags.Args()
		code := cmd.Run(cmd, args)
		os.Exit(code)
	}

	fmt.Fprintf(os.Stderr, "%s: unknown command %#q\n", a.name, args[0])
	fmt.Fprintf(os.Stderr, "Run `%s help` for usage.\n", a.name)
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
