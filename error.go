package cli

import (
	"fmt"
	"os"
)

// ErrFlagSyntax represents an error for bad arguments.
type ErrFlagSyntax string

// Error implements the error interface.
func (e ErrFlagSyntax) Error() string {
	return fmt.Sprintf("cli: invalid flag syntax '%s'", e)
}

// ErrUndefinedFlag represents an error for when an undefined flag is parsed.
type ErrUndefinedFlag string

// Error implements the error interface.
func (e ErrUndefinedFlag) Error() string {
	return fmt.Sprintf("cli: undefined flag '%s'", e)
}

// ErrRequiresArg represents an error for when an undefined flag is parsed.
type ErrRequiresArg string

// Error implements the error interface.
func (e ErrRequiresArg) Error() string {
	return fmt.Sprintf("cli: flag '%s' requires argument", e)
}

// Resolver represents the ability to resolve an error.
type Resolver func(err error)

// defaultResolver is the default resolver implementation.
func defaultResolver(err error) {
	if err != nil {
		os.Exit(1)
	}
}
