package cli

import (
	"html/template"
	"io"
	"strings"
)

// Usage represents application usage information.
type Usage struct {
	Name     string
	Usage    string
	Flags    []FlagUsage
	Commands []CommandUsage
}

// HasFlags returns true if global flags are available.
func (u Usage) HasFlags() bool {
	return len(u.Flags) > 0
}

// HasCommands returns true if global commands are available.
func (u Usage) HasCommands() bool {
	return len(u.Commands) > 0
}

// FlagUsage represents flag usage information.
type FlagUsage struct {
	Name    string
	Alias   string
	Usage   string
	Value   string
	Default string
}

// newFlagUsage returns the usage information for f.
func newFlagUsage(f *Flag) FlagUsage {
	return FlagUsage{
		Name:    "-" + f.name,
		Alias:   f.alias,
		Usage:   f.usage,
		Value:   f.value,
		Default: f.defaultValue,
	}
}

// CommandUsage represents command usage information.
type CommandUsage struct {
	Name  string
	Alias string
	Usage string
	Flags []FlagUsage
}

// HasFlags returns true if command flags are available.
func (u CommandUsage) HasFlags() bool {
	return len(u.Flags) > 0
}

// Summary returns the first line of the command usage information.
func (u CommandUsage) Summary() string {
	i := strings.Index(u.Usage, "\n")
	if i == -1 {
		return u.Usage
	}
	return u.Usage[:i]
}

// newCommandUsage returns the usage information for cmd.
func newCommandUsage(cmd *Command) CommandUsage {
	u := CommandUsage{
		Name:  cmd.name,
		Alias: cmd.alias,
		Usage: cmd.usage,
		Flags: make([]FlagUsage, len(cmd.flags)),
	}
	for i, f := range cmd.flags {
		u.Flags[i] = newFlagUsage(f)
	}
	return u
}

// UsageFormatter represents the ability to render usage information.
type UsageFormatter func(w io.Writer, u Usage) error

// defaultUsageFormatter is the default usage formatter implementation.
func defaultUsageFormatter(w io.Writer, u Usage) error {
	return tmpl(w, tmplUsage, u)
}

// tmpl parses text and applies data to it writing the output to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("tmpl")
	_, err := t.Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// tmplUsage represents the default application usage information template.
var tmplUsage = `{{.Usage}}

Usage:

    {{.Name}} [options] [command] [args...]

{{- if .HasFlags }}

Options:
{{range .Flags}}
    {{.Name | printf "%-16s"}} {{.Usage}}{{end}}
{{- end -}}

{{- if .HasCommands }}

Commands:
{{range .Commands}}
    {{.Name | printf "%-16s"}} {{.Summary}}{{end}}
{{- end }}

Run '{{.Name}} help [command]' for more information about a command.

`

// tmplCommandUsage represents the default command usage information template.
var tmplCommandUsage = `{{.Usage}}

{{ if .HasFlags -}}
Options:
{{range .Flags}}
    {{.Name | printf "%-16s"}} {{.Usage}}{{end}}
{{- end }}

`
