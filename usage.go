package cli

import (
	"errors"
	"io"
	"io/fs"
)

// nilUsage represents the nil usage.
type nilUsage struct{}

// Open implements the io/fs.FS interface.
func (u *nilUsage) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// UsageFS is a io/fs.FS implementation that
// reads files from usage lookup keys.
type UsageFS struct {
	fs    fs.FS
	ext   string
	index string
}

// NewUsageFS returns a usage lookup fs.FS implementation.
func NewUsageFS(fs fs.FS, opts ...UsageOption) fs.FS {
	if fs == nil {
		return &nilUsage{}
	}
	u := &UsageFS{fs: fs}
	for _, option := range opts {
		option(u)
	}
	if u.ext == "" {
		u.ext = ".md"
	}
	if u.index == "" {
		u.index = "README"
	}
	return u
}

// Open implements the io/fs.FS interface.
func (u *UsageFS) Open(name string) (fs.File, error) {
	if name == "" || name[len(name)-1] == '/' {
		name += u.index + u.ext
	} else {
		name += u.ext
	}
	return u.fs.Open(name)
}

// Usage displays the application usage information.
//
// The usage FS will be called with the help topic. The
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
	b, err := fs.ReadFile(c.usage, key)
	if err != nil {
		var perr *fs.PathError
		if !errors.As(err, &perr) {
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
