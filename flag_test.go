package cli

import "testing"

func TestNewFlag(t *testing.T) {
	var flag string
	f := NewFlag("flag", "test", &flag)
	if f == nil {
		t.Fatal("should return a flag")
	}
	if f.IsSet() {
		t.Fatal("should not be set")
	}
}

func TestNewFlagPanic(t *testing.T) {
	var flag string
	defer func() {
		perr := recover()
		if perr == nil {
			t.Fatal("should panic")
		}
	}()
	_ = NewFlag("flag", "test", flag)
}

func TestFlagSet(t *testing.T) {
	var flag string
	f := NewFlag("flag", "test", &flag)
	f.Set("test")
	if f.String() != "test" {
		t.Fatal("should set flag value")
	}
	if !f.IsSet() {
		t.Fatal("should be set")
	}
}

func TestFlagSetCount(t *testing.T) {
	var v bool
	f := NewFlag("verbose", "enable verbose output", &v, Bool(), ShortFlag("v"))
	f.Set("true")
	f.Set("true")
	f.Set("true")
	if f.Count() != 3 {
		t.Fatal("should increment set count")
	}
}
