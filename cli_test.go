package cli

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

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
			NewFlag("gs1", "global string 1", &c.gs1),
			NewFlag("gs2", "global string 2", &c.gs2),
			NewFlag("gb1", "global bool 1", &c.gb1, Bool()),
		}
		have, err := Parse(args, flags, true)
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
			NewFlag("gs1", "global string 1", &c.gs1),
			NewFlag("gs2", "global string 2", &c.gs2),
			NewFlag("gb1", "global bool 1", &c.gb1, Bool()),
		}
		t.Run("strict=on", func(t *testing.T) {
			want := ErrUndefinedFlag("undefined")
			_, err := Parse(args, flags, true)
			if !reflect.DeepEqual(err, want) {
				t.Fatalf("error for '%s'\nhave %v\nwant %v", line, err, want)
			}
		})
		t.Run("strict=off", func(t *testing.T) {
			_, err := Parse(args, flags, false)
			if err != nil {
				t.Fatalf("error for '%s'\nhave %v\nwant %v", line, err, nil)
			}
		})
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
			NewFlag("gs1", "global string 1", &c.gs1),
			NewFlag("gs2", "global string 2", &c.gs2),
			NewFlag("gb1", "global bool 1", &c.gb1, Bool()),
		}
		_, err := Parse(args, flags, true)
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
	flags := []*Flag{NewFlag("gs1", "global string 1", &c.gs1, EnvironmentKey(env))}
	_, err := Parse(args, flags, true)
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
	app := New("appname", "usage", nil)
	app.Add("foo", nil, "usage", nil)
}

type testCLI struct {
	gs1  string
	gs2  string
	gb1  bool
	args []string
}

func (c *testCLI) testCommandSuccess(args []string) error {
	return nil
}

var errCommandFailure = errors.New("cli: command failure")

func (c *testCLI) testCommandFailure(args []string) error {
	return errCommandFailure
}

func TestAddDuplicateCommand(t *testing.T) {
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatalf("duplicate command should panic")
		}
	}()
	c := &testCLI{}
	app := New("appname", "usage", nil)
	app.Add("test", c.testCommandSuccess, "usage", nil)
	app.Add("test", c.testCommandSuccess, "usage", nil)
}

func TestAddDuplicateCommandAlias(t *testing.T) {
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatalf("duplicate command alias should panic")
		}
	}()
	c := &testCLI{}
	app := New("appname", "usage", nil)
	app.Add("test", c.testCommandSuccess, "usage", nil)
	app.Add("aliased", c.testCommandSuccess, "usage", nil, Alias("test"))
}

func TestRunDefaultCommand(t *testing.T) {
	var testRunError bool
	opts := []Option{
		Stderr(ioutil.Discard),
		ErrorResolver(func(err error) { testRunError = true }),
	}
	app := New("appname", "usage", nil, opts...)
	app.Run([]string{"appname"})
	if !testRunError {
		t.Fatalf("default command should error")
	}
}

func TestRunCommandError(t *testing.T) {
	var testRunError error
	opts := []Option{
		Stderr(ioutil.Discard),
		ErrorResolver(func(err error) { testRunError = err }),
	}
	c := &testCLI{}
	app := New("appname", "usage", nil, opts...)
	app.Add("test", c.testCommandFailure, "usage", nil)
	app.Run([]string{"appname", "test"})
	if testRunError != errCommandFailure {
		t.Fatalf("Run error\nhave %v\nwant %v", testRunError, errCommandFailure)
	}
}
