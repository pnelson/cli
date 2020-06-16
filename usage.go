package cli

import (
	"fmt"
	"io"

	"github.com/charmbracelet/glamour"
)

// Renderer represents the ability to
// render help topics to terminal output.
type Renderer interface {
	Render(name string) ([]byte, error)
}

// defaultRenderer is the default Renderer implementation.
type defaultRenderer struct {
	data     map[string][]byte
	renderer *glamour.TermRenderer
}

// NewRenderer returns the default Renderer implementation.
//
// Usage information is rendered through a Markdown to ANSI
// compatible terminal renderer.
//
// The included cli-usage-gen program can generate the data
// from a directory of Markdown files, or you can roll your
// own per the Usage topic lookup convention.
func NewRenderer(data map[string][]byte) Renderer {
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		panic(fmt.Errorf("cli: failure to initialize renderer"))
	}
	return &defaultRenderer{
		data:     data,
		renderer: renderer,
	}
}

// Render implements the Renderer interface.
func (r *defaultRenderer) Render(name string) ([]byte, error) {
	b, ok := r.data[name]
	if !ok {
		return nil, ErrUsageNotFound
	}
	return r.renderer.RenderBytes(b)
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
