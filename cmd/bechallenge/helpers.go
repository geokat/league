package main

import (
	"strconv"
	"strings"
)

// Converts the slice of int literals `in` to an int slice.  Trims
// extraneous whitespace.
func atoi(in []string) ([]int, error) {
	out := make([]int, len(in))

	for i, s := range in {
		d, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}

		out[i] = d
	}

	return out, nil
}

// Converts the slice of ints `in` to a string of concatenated int literals.
func itos(in []int) string {
	var out string

	for _, d := range in {
		out += strconv.Itoa(d) + ","
	}

	if len(out) > 0 {
		// Get rid of the last comma.
		out = out[:len(out)-1]
	}

	return out
}
