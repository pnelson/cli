# cli

Package cli provides structure for command line applications with sub-commands.

This package was heavily influenced by the Go command line application.
Commands help and version are implemented by default. The usage information
is pretty printed in an opinionated format.

The most basic cli application is boring:

```go
app := cli.New("app", "0.0.1")
app.Run()
```  

Build and run your application:

```
$ ./app
Usage: app <command> [options] [<args>]

    help        Output this usage information.
    version     Output the application version.

Use "app help [command]" for more information about a command.
```

Add some commands:

```go
app.Command(&cli.Command{
  Usage: "say <message>",
  Short: "Print a message to stdout.",
  Run: func(args []string) int {
    if len(args) < 1 {
      return 1
    }
    fmt.Println(args[0])
    return 0
  },
})
```

Let's build and run the application again:

```
$ ./app
Usage: app <command> [options] [<args>]

    help        Output this usage information.
    version     Output the application version.
    say         Print a message to stdout.

Use "app help [command]" for more information about a command.
```

Awesome. The say command seems useful. Let's try it:

```
$ ./app say 'hello, world'
hello, world
```