package cli

import (
	"bytes"
	"embed"
	"io/fs"
	"io/ioutil"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func newTestUsage(t *testing.T) fs.FS {
	t.Helper()
	usage, err := fs.Sub(testdata, "testdata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return NewUsageFS(usage)
}

func TestUsage(t *testing.T) {
	var tests = []struct {
		scope string
		args  []string
		want  string
	}{
		{
			"",
			[]string{"appname", "help"},
			"README.md\n",
		},
		{
			"",
			[]string{"appname", "help", "test"},
			"test.md\n",
		},
		{
			"cli",
			[]string{"appname", "help"},
			"cli/README.md\n",
		},
		{
			"cli",
			[]string{"appname", "help", "test"},
			"cli/test.md\n",
		},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		opts := []Option{Scope(tt.scope), Stdout(&buf), Stderr(ioutil.Discard)}
		app := New("appname", newTestUsage(t), nil, opts...)
		app.Add("test", testCommand, nil)
		err := app.Run(tt.args)
		if err != nil {
			t.Fatalf("help args=%v scope='%s'\ncommand should not error", tt.args, tt.scope)
		}
		have := buf.String()
		if have != tt.want {
			t.Fatalf("help %v\nhave '%s'\nwant '%s'", tt.args, have, tt.want)
		}
	}
}

func TestUsageNil(t *testing.T) {
	app := New("appname", nil, nil, Stdout(ioutil.Discard), Stderr(ioutil.Discard))
	err := app.Usage(nil, "test")
	if err != ErrExitFailure {
		t.Fatalf("usage should error")
	}
}

func TestUsageRoot(t *testing.T) {
	for _, scope := range []string{"", "cli"} {
		var buf bytes.Buffer
		opts := []Option{Scope(scope), Stdout(ioutil.Discard), Stderr(&buf)}
		app := New("appname", newTestUsage(t), nil, opts...)
		err := app.Run([]string{"appname"})
		if err != ErrExitFailure {
			t.Fatalf("help root command scope='%s'\ncommand should error", scope)
		}
		have := buf.String()
		want := "README.md\n"
		if have != want {
			t.Fatalf("help root command scope='%s'\nhave '%s'\nwant '%s'", scope, have, want)
		}
	}
}

func TestUsageNotFound(t *testing.T) {
	app := New("appname", newTestUsage(t), nil, Stdout(ioutil.Discard), Stderr(ioutil.Discard))
	app.Add("test", func([]string) error { return nil }, nil)
	err := app.Run([]string{"appname", "help", "not-found"})
	if err != ErrExitFailure {
		t.Fatalf("help command should error")
	}
}
