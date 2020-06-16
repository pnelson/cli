package cli

import (
	"reflect"
	"strings"
)

// Flag represents a flag.
type Flag struct {
	flag         reflect.Value
	kind         FlagKind
	name         string
	alias        string
	count        int
	value        string
	envKey       string
	defaultValue string
}

// NewFlag returns a new flag. The flag must be a pointer. You must pass the
// Kind option so the flag parser knows how to process the command line unless
// the flag points to a string, the default flag kind.
func NewFlag(name string, flag interface{}, opts ...FlagOption) *Flag {
	v := reflect.ValueOf(flag)
	if v.Kind() != reflect.Ptr {
		panic("cli: flag must be pointer")
	}
	f := &Flag{
		flag: v.Elem(),
		kind: flagString{},
		name: strings.ToLower(name),
	}
	for _, option := range opts {
		option(f)
	}
	return f
}

// Count returns the number of times the flag was set.
func (f *Flag) Count() int {
	return f.count
}

// IsSet returns true if the flag was explicitly set.
func (f *Flag) IsSet() bool {
	return f.count > 0
}

// Set sets the flag value.
func (f *Flag) Set(value string) {
	f.count++
	f.flag.Set(reflect.ValueOf(f.kind.Parse(value)))
	f.value = value
}

// String returns the value as a string. Boolean flags are
// returned as "true" or "false" as strconv.FormatBool would.
//
// String implements the fmt.Stringer interface.
func (f *Flag) String() string {
	if f == nil {
		return ""
	}
	return f.value
}

// FlagKind represents the type of flag.
type FlagKind interface {
	Parse(value string) interface{}
	HasArg() bool
}

// flagString represents a string flag.
type flagString struct{}

// Parse returns the value as-is.
//
// Parse implements the FlagKind interface.
func (f flagString) Parse(value string) interface{} {
	return value
}

// HasArg implements the FlagKind interface.
func (f flagString) HasArg() bool {
	return true
}

// flagBool represents a boolean flag.
type flagBool struct{}

// Parse returns "true" if the value is
// 1, t, T, true, TRUE, True, y, Y, yes, YES, Yes.
//
// Parse implements the FlagKind interface.
func (f flagBool) Parse(value string) interface{} {
	switch value {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes":
		return true
	}
	return false
}

// HasArg implements the FlagKind interface.
func (f flagBool) HasArg() bool {
	return false
}

// FlagOption represents a functional option for flag configuration.
type FlagOption func(*Flag)

// Kind sets the flag kind.
// This option is required unless the flag points to a string.
// This option must be used before the DefaultValue option.
func Kind(kind FlagKind) FlagOption {
	return func(f *Flag) {
		f.kind = kind
	}
}

// Bool sets the flag kind to the built in boolean flag kind.
func Bool() FlagOption {
	return Kind(flagBool{})
}

// ShortFlag sets the short flag.
func ShortFlag(name string) FlagOption {
	return func(f *Flag) {
		f.alias = name
	}
}

// DefaultValue sets the flag default value.
func DefaultValue(value string) FlagOption {
	return func(f *Flag) {
		f.flag.Set(reflect.ValueOf(f.kind.Parse(value)))
		f.value = value
		f.defaultValue = value
	}
}

// EnvironmentKey sets the flag environment variable key.
func EnvironmentKey(key string) FlagOption {
	return func(f *Flag) {
		f.envKey = key
	}
}
