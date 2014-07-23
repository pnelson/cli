package cli

import (
	"io"
	"strings"
	"text/template"
)

var usageTemplate = `Usage: {{.Name}} <command> [options] [<args>]
{{range .Commands}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "{{.Name}} help [command]" for more information about a command.

`

var helpTemplate = `usage: {{.Name}} {{.Command.Usage}}{{if .Command.Long}}

{{.Command.Long | trim}}{{end}}
`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("usage")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})

	template.Must(t.Parse(text))

	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
