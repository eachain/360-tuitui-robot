package main

import (
	"fmt"
	"testing"
)

func equalStrings(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("expect size: %v, return size: %v", len(a), len(b))
	}

	for i := range a {
		if a[i] != b[i] {
			return fmt.Errorf("%vth expect: %q, return: %q", i, a[i], b[i])
		}
	}
	return nil
}

func TestParseTokens(t *testing.T) {
	var cases = []struct {
		S string
		T []string
	}{{
		``,
		[]string{},
	}, {
		`a`,
		[]string{"a"},
	}, {
		`   a   `,
		[]string{"a"},
	}, {
		`a b c`,
		[]string{"a", "b", "c"},
	}, {
		`a -b=123 -c`,
		[]string{"a", "-b=123", "-c"},
	}, {
		`"1" '2' 3\n4 "5\t6" '7\n8' 9\ 0 a'b'c "d'e'f" g"h"i 'j"k"l'`,
		[]string{"1", "2", "3\n4", "5\t6", `7\n8`, "9 0", "abc", "d'e'f", `ghi`, `j"k"l`},
	}}

	for _, c := range cases {
		tokens, err := ParseTokens(c.S)
		if err != nil {
			t.Fatalf("parse %q: %v", c.S, err)
		}
		if err = equalStrings(c.T, tokens); err != nil {
			t.Fatalf("parse %q: %v", c.S, err)
		}
	}
}
