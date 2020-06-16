package cli

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

var testUsage = testRenderer{
	"":         []byte("README.md"),
	"test":     []byte("test.md"),
	"cli/":     []byte("cli/README.md"),
	"cli/test": []byte("cli/test.md"),
}

type testRenderer map[string][]byte

func (u testRenderer) Render(name string) ([]byte, error) {
	v, ok := u[name]
	if !ok {
		return nil, ErrUsageNotFound
	}
	return v, nil
}

func TestUsage(t *testing.T) {
	var tests = []struct {
		scope string
		args  []string
		want  []byte
	}{
		{
			"",
			[]string{"appname", "help"},
			testUsage[""],
		},
		{
			"",
			[]string{"appname", "help", "test"},
			testUsage["test"],
		},
		{
			"cli",
			[]string{"appname", "help"},
			testUsage["cli/"],
		},
		{
			"cli",
			[]string{"appname", "help", "test"},
			testUsage["cli/test"],
		},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		opts := []Option{UsageScope(tt.scope), Stdout(&buf), Stderr(ioutil.Discard)}
		app := New("appname", testUsage, nil, opts...)
		app.Add("test", testCommand, nil)
		err := app.Run(tt.args)
		if err != nil {
			t.Fatalf("help args=%v scope='%s'\ncommand should not error", tt.args, tt.scope)
		}
		have := buf.Bytes()
		if !reflect.DeepEqual(have, tt.want) {
			t.Fatalf("help %v\nhave '%s'\nwant '%s'", tt.args, have, tt.want)
		}
	}
}

func TestUsageRoot(t *testing.T) {
	for _, scope := range []string{"", "cli"} {
		var buf bytes.Buffer
		opts := []Option{UsageScope(scope), Stdout(ioutil.Discard), Stderr(&buf)}
		app := New("appname", testUsage, nil, opts...)
		err := app.Run([]string{"appname"})
		if err != ErrExitFailure {
			t.Fatalf("help root command scope='%s'\ncommand should error", scope)
		}
		have := buf.Bytes()
		want := testUsage[""]
		if !reflect.DeepEqual(have, want) {
			t.Fatalf("help root command scope='%s'\nhave '%s'\nwant '%s'", scope, have, want)
		}
	}
}

func TestUsageNotFound(t *testing.T) {
	app := New("appname", testUsage, nil, Stdout(ioutil.Discard), Stderr(ioutil.Discard))
	app.Add("test", func([]string) error { return nil }, nil)
	err := app.Run([]string{"appname", "help", "not-found"})
	if err != ErrExitFailure {
		t.Fatalf("help command should error")
	}
}
