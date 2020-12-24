package cli

// Command represents an application command.
type Command struct {
	name       string
	alias      string
	proxy      bool
	flags      []*Flag
	handler    Handler
	middleware []func(Handler) Handler
}

// Handler represents a command handler.
type Handler func(args []string) error

// NewCommand returns a new command.
func NewCommand(name string, handler Handler, flags []*Flag, opts ...CommandOption) *Command {
	c := &Command{
		name:       name,
		flags:      flags,
		middleware: make([]func(Handler) Handler, 0),
	}
	for _, option := range opts {
		option(c)
	}
	c.build(handler)
	return c
}

// build wraps h with the configured middleware.
func (c *Command) build(h Handler) {
	c.handler = h
	for i := len(c.middleware) - 1; i >= 0; i-- {
		c.handler = c.middleware[i](c.handler)
	}
}

// CommandOption represents a functional option for command configuration.
type CommandOption func(*Command)

// Alias sets the command alias.
func Alias(name string) CommandOption {
	return func(c *Command) {
		c.alias = name
	}
}

// Proxy instructs the dispatcher to proxy the unparsed
// arguments to the command itself for further processing.
func Proxy() CommandOption {
	return func(c *Command) {
		c.proxy = true
	}
}

// WithMiddleware appends middleware to the middleware stack.
func WithMiddleware(middleware ...func(Handler) Handler) CommandOption {
	return func(c *Command) {
		c.middleware = append(c.middleware, middleware...)
	}
}
