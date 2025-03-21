package main

import (
	"fmt"
	"math/big"
	"strings"
)

// Converts the slice of int literals `in` to an big.Int slice.  Trims
// extraneous whitespace.
func atoi(in []string) ([]*big.Int, error) {
	out := make([]*big.Int, len(in))

	for i, s := range in {
		s = strings.TrimSpace(s)
		d := new(big.Int)

		if _, ok := d.SetString(s, 10); !ok {
			return nil, fmt.Errorf(`parsing "%s": invalid syntax`, s)
		}

		out[i] = d
	}

	return out, nil
}

// Converts the slice of big.Int's `in` to a string of concatenated
// int literals.
func itos(in []*big.Int) string {
	var out string

	for _, d := range in {
		out += d.String() + ","
	}

	if len(out) > 0 {
		// Get rid of the last comma.
		out = out[:len(out)-1]
	}

	return out
}
