package cli

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

type testCLI struct {
	gs1  string
	gs2  string
	gb1  bool
	args []string
}

var errCommandFailure = errors.New("cli: command failure")

func testCommand(args []string) error {
	return nil
}

func testCommandFailure(args []string) error {
	return errCommandFailure
}

func testCommandErrUsage(args []string) error {
	return ErrUsage
}

func TestParse(t *testing.T) {
	tests := map[string]*testCLI{
		"":                   &testCLI{},
		"-gs1 string":        &testCLI{gs1: "string"},
		"-gs1=string":        &testCLI{gs1: "string"},
		"-gb1":               &testCLI{gb1: true},
		"-gs1 string -gb1":   &testCLI{gs1: "string", gb1: true},
		"-gb1 -gs1 string":   &testCLI{gs1: "string", gb1: true},
		"--gb1 --gs1 string": &testCLI{gs1: "string", gb1: true},
	}
	for line, want := range tests {
		c := &testCLI{}
		args := strings.Split(line, " ")
		flags := []*Flag{
			NewFlag("gs1", &c.gs1),
			NewFlag("gs2", &c.gs2),
			NewFlag("gb1", &c.gb1, Bool()),
		}
		have, err := Parse(args, flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(c, want) {
			t.Fatalf("flags for '%s'\nhave %v\nwant %v", line, c, want)
		}
		if !reflect.DeepEqual(have, c.args) {
			if c.args != nil {
				t.Fatalf("args for '%s'\nhave %v\nwant %v", line, have, c.args)
			}
		}
	}
}

func TestParseArgs(t *testing.T) {
	tests := map[string][]string{
		"-gs1 string":        []string{},
		"-gs1=string":        []string{},
		"-gb1":               []string{},
		"-gs1 string -gb1":   []string{},
		"-gb1 -gs1 string":   []string{},
		"--gb1 --gs1 string": []string{},
		"-":                  []string{"-"},
		"--":                 []string{"--"},
		"-- arg":             []string{"--", "arg"},
		"arg":                []string{"arg"},
		"-gs1 string -":      []string{"-"},
		"-gs1 string --":     []string{"--"},
		"-gs1 string -- arg": []string{"--", "arg"},
		"-gs1 string arg":    []string{"arg"},
	}
	for line, want := range tests {
		c := &testCLI{}
		args := strings.Split(line, " ")
		flags := []*Flag{
			NewFlag("gs1", &c.gs1),
			NewFlag("gs2", &c.gs2),
			NewFlag("gb1", &c.gb1, Bool()),
		}
		have, err := Parse(args, flags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(have, want) {
			t.Fatalf("args for '%s'\nhave %v\nwant %v", line, have, want)
		}
	}
}

func TestParseUndefined(t *testing.T) {
	tests := map[string]struct{}{
		"-undefined":                    struct{}{},
		"-undefined string":             struct{}{},
		"-undefined=string":             struct{}{},
		"-gs1 string -undefined":        struct{}{},
		"-gs1 string -undefined string": struct{}{},
		"-gs1 string -undefined=string": struct{}{},
		"-undefined -gs1 string":        struct{}{},
		"-undefined string -gs1 string": struct{}{},
		"-undefined=string -gs1 string": struct{}{},
	}
	for line := range tests {
		c := &testCLI{}
		args := strings.Split(line, " ")
		flags := []*Flag{
			NewFlag("gs1", &c.gs1),
			NewFlag("gs2", &c.gs2),
			NewFlag("gb1", &c.gb1, Bool()),
		}
		want := ErrUndefinedFlag("undefined")
		_, err := Parse(args, flags)
		if !reflect.DeepEqual(err, want) {
			t.Fatalf("should return undefined flag error for '%s'", line)
		}
	}
}

func TestParseRequiresArg(t *testing.T) {
	tests := map[string]struct{}{
		"-gs1":        struct{}{},
		"-gs1 -gb1":   struct{}{},
		"-gb1 -gs1":   struct{}{},
		"--gs1":       struct{}{},
		"--gs1 --gb1": struct{}{},
		"--gb1 --gs1": struct{}{},
	}
	want := ErrRequiresArg("gs1")
	for line := range tests {
		c := &testCLI{}
		args := strings.Split(line, " ")
		flags := []*Flag{
			NewFlag("gs1", &c.gs1),
			NewFlag("gs2", &c.gs2),
			NewFlag("gb1", &c.gb1, Bool()),
		}
		_, err := Parse(args, flags)
		if !reflect.DeepEqual(err, want) {
			t.Fatalf("error for '%s'\nhave %v\nwant %v", line, err, want)
		}
	}
}

func TestParseEnv(t *testing.T) {
	const env = "TEST_PARSE_ENV"
	const want = "env"
	os.Setenv(env, want)
	c := &testCLI{}
	args := []string{}
	flags := []*Flag{NewFlag("gs1", &c.gs1, EnvironmentKey(env))}
	_, err := Parse(args, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.gs1 != want {
		t.Fatalf("gs1\nhave '%s'\nwant '%s'", c.gs1, want)
	}
}

func TestAddNilHandler(t *testing.T) {
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatalf("nil handler should panic")
		}
	}()
	app := New("appname", testUsage, nil)
	app.Add("foo", nil, nil)
}

func TestAddDuplicateCommand(t *testing.T) {
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatalf("duplicate command should panic")
		}
	}()
	app := New("appname", testUsage, nil)
	app.Add("test", testCommand, nil)
	app.Add("test", testCommand, nil)
}

func TestAddDuplicateCommandAlias(t *testing.T) {
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatalf("duplicate command alias should panic")
		}
	}()
	app := New("appname", testUsage, nil)
	app.Add("test", testCommand, nil)
	app.Add("aliased", testCommand, nil, Alias("test"))
}

func TestRunDefaultCommand(t *testing.T) {
	var buf bytes.Buffer
	app := New("appname", testUsage, nil, Stderr(&buf))
	err := app.Run([]string{"appname"})
	if err != ErrExitFailure {
		t.Fatalf("default command should error")
	}
	have := buf.Bytes()
	want := testUsage[""]
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("should return root usage docs\nhave '%s'\nwant '%s'", have, want)
	}
}

func TestRunHelpError(t *testing.T) {
	app := New("appname", testUsage, nil, Stderr(ioutil.Discard))
	app.Add("test", func([]string) error { return nil }, nil)
	err := app.Run([]string{"appname", "help", "test", "fail"})
	if err != ErrExitFailure {
		t.Fatalf("help command should error")
	}
}

func TestRunCommandError(t *testing.T) {
	app := New("appname", testUsage, nil, Stderr(ioutil.Discard))
	app.Add("test", testCommandFailure, nil)
	err := app.Run([]string{"appname", "test"})
	if err != errCommandFailure {
		t.Fatalf("Run error\nhave %v\nwant %v", err, errCommandFailure)
	}
}

func TestRunCommandErrUsage(t *testing.T) {
	var buf bytes.Buffer
	c := &testCLI{}
	flags := []*Flag{
		NewFlag("gs1", &c.gs1),
		NewFlag("gs2", &c.gs2),
		NewFlag("gb1", &c.gb1, Bool()),
	}
	app := New("appname", testUsage, flags, Stderr(&buf))
	app.Add("test", testCommandErrUsage, nil)
	err := app.Run([]string{"appname", "-gb1", "-gs1", "string", "test"})
	if err != ErrExitFailure {
		t.Fatalf("Run error\nhave %v\nwant %v", err, ErrExitFailure)
	}
	have := buf.Bytes()
	want := testUsage["test"]
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("should return test command usage docs\nhave '%s'\nwant '%s'", have, want)
	}
}
