package cli

import "io"

// Option represents a functional option for configuration.
type Option func(*CLI)

// Help sets the application help handler.
func Help(handler Handler) Option {
	return func(c *CLI) {
		c.helpHandler = handler
	}
}

// Version enables the application version handler.
func Version(version string) Option {
	return func(c *CLI) {
		c.version = version
	}
}

// Scope sets the help topic scope for registered commands.
// See Usage documentation for more information.
func Scope(scope string) Option {
	return func(c *CLI) {
		c.scope = scope
	}
}

// Default sets the handler to execute when no command is given.
func Default(handler Handler) Option {
	return func(c *CLI) {
		c.defaultHandler = handler
	}
}

// Stdin sets the stdin reader. Defaults to os.Stdin.
// A nil reader will fallback to os.Stdin.
func Stdin(r io.Reader) Option {
	return func(c *CLI) {
		c.stdin = r
	}
}

// Stdout sets the stdout writer. Defaults to os.Stdout.
// A nil writer will fallback to os.Stdout.
// Use ioutil.Discard to discard output.
func Stdout(w io.Writer) Option {
	return func(c *CLI) {
		c.stdout = w
	}
}

// Stderr sets the stderr writer. Defaults to os.Stderr.
// A nil writer will fallback to os.Stderr.
// Use ioutil.Discard to discard output.
func Stderr(w io.Writer) Option {
	return func(c *CLI) {
		c.stderr = w
	}
}
