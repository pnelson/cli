# cli

Package cli provides structure for command line applications with commands.

## About

This package is highly opinionated but I have tried to provide hooks in some
key areas, though I haven't completely thought out if it is sufficient. The
package is built using variadic options so it is easy to extend.

One thing I really wanted was to keep the command definitions, usage
information, and flags all together. I think it's easier to maintain for small
to medium sized applications and it is still easy to break up if the
application outgrows that pattern.

I've opted for the simplicity of only a single level of commands. This may
change in the future but I'm really just trying to build for the 99% here. For
now, there are enough packages in the ecosystem that support nested commands.

My design was based around having a readable and maintainable way to:

- define commands
- define global (application) flags
- define local (command) flags
- fallback to prefixed environment variables
- output formatted usage

Custom flag types are easy to implement but I felt the need to depart from the
standard library `flag` package interfaces to provide the developer experience
I'm going for. See the `FlagKind` interface for details.

Generating shell completion and manpages should come eventually.

## Usage

The most basic cli application is boring so we'll take it a step further and
add a global flag with a default value. Global flags are applicable to any
command and must come before the command name on the command line. Single and
double flags are supported.

```go
var user string
app := cli.New("appname", `Usage summary recommended 50 character max length

Include a more detailed description, if necessary. Set your editor
to wrap lines at 72 characters. This text will be trimmed of spaces.
`, []*cli.Flag{
  cli.NewFlag("user", "user", &user, cli.DefaultValue("demo")),
}, cli.Version("0.0.1"))
app.Run(os.Args)
```  

You can get creative and put the flags in a struct and write the handlers as
methods on the struct to gain access to the flags.

Build and run your application:

```
$ ./appname
Usage summary recommended 50 character max length

Include a more detailed description, if necessary. Set your editor
to wrap lines at 72 characters. This text will be trimmed of spaces.

Usage:

    appname [options] [command] [args...]

Options:

    -user       user

Commands:

    help        show command usage information
    version     show version information

Run 'appname help [command]' for more information about a command.

$ echo $?
1
```

The default command prints the usage information and returns as an error.
Override this behavior with the `Default` option.

You'll notice two commands added by default there. Well, `help` was default.
We added the version command when we specified the `Version` option.

Let's add a command `say` that prints the first argument if it exists. For
completeness, we'll also add two flags, boolean since they don't have args.
This command is a little like `echo` so we'll use the `Alias` option.

```go
func say(args []string) error {
  if len(args) < 1 {
    return errors.New("say requires arguments")
  }
  line := strings.Join(args, " ")
  if upper {
    line = strings.ToUpper(line)
  }
  if username != "" {
    fmt.Printf("%s: ", username)
  }
  fmt.Printf("%s", line)
  if !newline {
    fmt.Printf("\n")
  }
  return nil
}

var upper bool
var newline bool
app.Add("say", say, "display a line of text", []*cli.Flag{
  cli.NewFlag("upper", "transform characters to upper case", &upper, cli.Bool()),
  cli.NewFlag("newline", "do not output the trailing newline", &newline, cli.Bool(), cli.ShortFlag("n")),
}, cli.Alias("echo"))
```

The error return value of the dispatched handler is passed to the error
resolver which by default calls `os.Exit(1)` if the error is non-`nil`.
Override this behavior with the `ErrorResolver` option.

Let's build and run the application again:

```
$ ./appname
Usage summary recommended 50 character max length

Include a more detailed description, if necessary. Set your editor
to wrap lines at 72 characters. This text will be trimmed of spaces.

Usage:

    appname [options] [command] [args...]

Options:

    -user       user

Commands:

    help        show command usage information
    say         display a line of text
    version     show version information

Run 'appname help [command]' for more information about a command.
```

I'm a bit curious what the `say` command is all about. Let's get some help:

```
$ ./appname help say
display a line of text

Options:

    -upper      transform characters to upper case
    -newline    do not output the trailing newline
```

I have a good idea of how the application works now. Let's try it:

```
$ appname version
0.0.1

$ appname say 'Hello, world'
demo: Hello, world

$ appname echo -upper 'Hello, world'
demo: HELLO, WORLD

$ APPNAME_UPPER="true" appname say 'Hello, world'
demo: HELLO, WORLD

$ APPNAME_USER="" APPNAME_UPPER="true" appname say 'Hello, world'
HELLO, WORLD

$ APPNAME_USER="" APPNAME_UPPER="true" appname -user 'nobody' say 'Hello, world'
nobody: HELLO, WORLD
```

Explicitly defined flag values on the command line take precendence over
environment variables and default values.
