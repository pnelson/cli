# cli

Package cli provides structure for command line applications.

This package is moderately opinionated but I have tried to provide hooks in some
key areas. The package is built using variadic options so it is easy to extend.

I wanted to keep the command and flag definitions together. I think it's easier
to maintain for small to medium sized applications and it is still easy to
break up if the application outgrows that pattern.

I value documentation so I've separated the usage lookup so as to not subtly
encourage minimal documentation. A nice side effect of this decision is that it
should be relatively straight forward to plug in internationalization support
for the built in help command.

I've opted for the simplicity of only a single level of commands. This may
change in the future but I'm really just trying to build for the 99% here. For
now, there are enough packages in the ecosystem that support nested commands.

My design was based around having a readable and maintainable way to:

- define commands
- define global (application) flags
- define local (command) flags
- create environment variable flag mappings
- fallback to application name prefixed environment variables
- output formatted usage
- testable

Custom flag types are easy to implement but I felt the need to depart from the
standard library `flag` package interfaces to provide the developer experience
I'm going for. See the `FlagKind` interface for details.

Explicitly defined flag values on the command line take precedence over
environment variables and default values.

Generating shell completion should come eventually.
