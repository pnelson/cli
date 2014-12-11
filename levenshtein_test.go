package cli

import "testing"

func TestLevenshtein(t *testing.T) {
	var tests = []struct {
		v int
		s string
		t string
	}{
		{3, "kitten", "sitting"},
		{3, "Saturday", "Sunday"},
	}

	for i, tt := range tests {
		v := levenshtein(tt.s, tt.t)
		if v != tt.v {
			t.Errorf("%d. levenshtein(%q, %q)\nhave %d\nwant %d", i, tt.s, tt.t, v, tt.v)
		}
	}
}
