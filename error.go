package cli

import "fmt"

// ErrExitFailure represents errors that should immediately
// exit with failure status. All output to stdout or stderr
// should be written before a handler returns this value, as
// is the case for built in usage. It is the package user's
// responsibility to handle this error.
var ErrExitFailure = fmt.Errorf("1")

// ErrUsage represents the error that should be returned
// by handlers to output usage information for the command.
// ErrUsage will be rewritten as ErrExitFailure on success.
var ErrUsage = fmt.Errorf("cli: usage")

// ErrFlagSyntax represents an error for bad arguments.
type ErrFlagSyntax string

// Error implements the error interface.
func (e ErrFlagSyntax) Error() string {
	return fmt.Sprintf("Flag '%s' is syntactically incorrect.", string(e))
}

// ErrUndefinedFlag represents an error for when an undefined flag is parsed.
type ErrUndefinedFlag string

// Error implements the error interface.
func (e ErrUndefinedFlag) Error() string {
	return fmt.Sprintf("Flag '%s' is undefined.", string(e))
}

// ErrRequiresArg represents an error for when an undefined flag is parsed.
type ErrRequiresArg string

// Error implements the error interface.
func (e ErrRequiresArg) Error() string {
	return fmt.Sprintf("Flag '%s' requires an argument.", string(e))
}
