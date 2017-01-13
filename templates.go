package cli

import (
	"io"
	"strings"
	"text/template"
)

type usageData struct {
	Name     string
	Commands []*Command
}

var usageTemplate = `Usage: {{.Name}} <command> [options] [<args>]
{{range .Commands}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "{{.Name}} help [command]" for more information about a command.

`

type helpData struct {
	Name    string
	Command *Command
}

var helpTemplate = `usage: {{.Name}} {{.Command.Usage}}{{if .Command.Long}}

{{.Command.Long | trim}}{{end}}
`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("tmpl")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})
	template.Must(t.Parse(text))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
