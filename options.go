package cli

import "fmt"

// Option represents a functional option for configuration.
type Option func(*Application)

// Version sets the version string and adds a version command.
func Version(version string) Option {
	return func(app *Application) {
		app.Command(&Command{
			native: true,
			Usage:  "version",
			Short:  "Output the application version.",
			Run: func(args []string) int {
				fmt.Printf("%s v%s\n", app.name, version)
				return 0
			},
		})
	}
}
