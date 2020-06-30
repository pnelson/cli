package cli

import "io"

// Renderer represents the ability to
// render help topics to terminal output.
type Renderer interface {
	Render(name string) ([]byte, error)
}

// defaultRenderer is the default Renderer implementation.
type defaultRenderer map[string][]byte

// NewRenderer returns the default Renderer implementation.
//
// Usage information is rendered through a Markdown to ANSI
// compatible terminal renderer.
//
// The included cli-usage-gen program can generate the data
// from a directory of Markdown files, or you can roll your
// own per the Usage topic lookup convention.
func NewRenderer(data map[string][]byte) Renderer {
	return defaultRenderer(data)
}

// Render implements the Renderer interface.
func (r defaultRenderer) Render(name string) ([]byte, error) {
	b, ok := r[name]
	if !ok {
		return nil, ErrUsageNotFound
	}
	return b, nil
}

// Usage displays the application usage information.
//
// The renderer will be called with the help topic. The
// help topic is prefixed with the configured scope if the
// topic is a registered command. For example, if the scope
// is "cli" and the "foo" command is registered, "help foo"
// will call the renderer with "cli/foo" but "help not-found"
// would passthrough as "not-found" without the scope.
func (c *CLI) Usage(w io.Writer, name string) error {
	key := name
	_, ok := c.commands[name]
	if ok {
		key = c.scope + name
	}
	b, err := c.usage.Render(key)
	if err != nil {
		if err != ErrUsageNotFound {
			return err
		}
		if name == "" || name == c.scope {
			c.Errorf("Unknown help topic.\n")
		} else {
			c.Errorf("Unknown help topic '%s'.\n", name)
		}
		c.Errorf("Run '%s help' for usage information.\n", c.name)
		return ErrExitFailure
	}
	_, err = w.Write(b)
	return err
}
