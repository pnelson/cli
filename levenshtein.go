package cli

// similarThreshold is the maximum levenshtein distance at
// which a command is considered to be similar.
const similarThreshold = 5

// levenshtein returns the levenshtein distance of s from t.
func levenshtein(s, t string) int {
	if s == t {
		return 0
	}
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}
	v0 := make([]int, len(t)+1)
	v1 := make([]int, len(t)+1)
	for i := 0; i < len(v0); i++ {
		v0[i] = i
	}
	for i := 0; i < len(s); i++ {
		v1[0] = i + 1
		for j := 0; j < len(t); j++ {
			cost := 0
			if s[i] != t[j] {
				cost = 1
			}
			v1[j+1] = min(v1[j]+1, v0[j+1]+1, v0[j]+cost)
		}
		for j := 0; j < len(v0); j++ {
			v0[j] = v1[j]
		}
	}
	return v1[len(t)]
}

// min returns the minimum of one or more integers.
func min(xs ...int) int {
	m := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] < m {
			m = xs[i]
		}
	}
	return m
}
