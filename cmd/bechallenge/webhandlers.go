package main

import (
	"fmt"
	h "net/http"
)

// Handles echo requests by validating the parsed matrix of int
// literals and returning it back. Expects the matrix CSV in the
// request context.
func handleEcho(w h.ResponseWriter, r *h.Request) {
	echo(w, r, false)
}

// Handles flatten requests by validating the parsed matrix of int
// literals and returning a one line string with their
// concatenation. Expects the matrix CSV in the request context.
func handleFlatten(w h.ResponseWriter, r *h.Request) {
	echo(w, r, true)
}

// Handles sum requests by validating the supplied matrix of int
// literals and returning a string with their sum. Expects the matrix
// CSV in the request context.
func handleSum(w h.ResponseWriter, r *h.Request) {
	reduce(w, r, false)
}

// Handles multiply requests by validating the supplied matrix of int
// literals and returning a string with their product. Expects the
// matrix CSV in the request context.
func handleMultiply(w h.ResponseWriter, r *h.Request) {
	reduce(w, r, true)
}

// Implements the actual handler for reduce-like (sum, multiply)
// requests.
func reduce(w h.ResponseWriter, r *h.Request, multiply bool) {
	recs := r.Context().Value(csvRecordsKey).([][]string)

	// Handle zero size matrix edge case.
	if len(recs) == 0 {
		fmt.Fprint(w, 0, "\n")

		return
	}

	var resp int

	// Zero value of int would turn every product to 0.
	if multiply {
		resp = 1
	}

	// Process the CSV rows and build the response inline.
	for ri, row := range recs {
		ints, err := atoi(row)
		if err != nil {
			h.Error(
				w,
				fmt.Sprintf("Error: parsing CSV: record on line %d: %v", ri+1, err),
				h.StatusBadRequest)

			return
		}

		for _, int := range ints {
			if multiply {
				resp *= int
			} else {
				resp += int
			}
		}
	}

	// The challenge description requires a trailing "\n" in the response.
	fmt.Fprint(w, resp, "\n")
}

// Implements the actual handler for echo-like (echo, flatten)
// requests.
func echo(w h.ResponseWriter, r *h.Request, flatten bool) {
	recs := r.Context().Value(csvRecordsKey).([][]string)

	var resp string

	// Process the CSV rows and build the response inline.
	for ri, row := range recs {
		ints, err := atoi(row)
		if err != nil {
			h.Error(
				w,
				fmt.Sprintf("Error: parsing CSV: record on line %d: %v", ri+1, err),
				h.StatusBadRequest)

			return
		}

		resp += itos(ints)
		if flatten {
			resp += ","
		} else {
			resp += "\n"
		}
	}

	if len(resp) == 0 {
		// The challenge description requires a trailing "\n" in the response.
		resp += "\n"
	} else if flatten {
		// Replace last comma with a new line
		resp = resp[:len(resp)-1] + "\n"
	}
	fmt.Fprint(w, resp)
}

// Handles invert requests by validating the supplied matrix of int
// literals and returning its transpose. Expects the matrix CSV in the
// request context.
func handleInvert(w h.ResponseWriter, r *h.Request) {
	recs := r.Context().Value(csvRecordsKey).([][]string)

	// Transposed matrix.
	tran := make([][]int, len(recs))

	// Process the CSV rows.
	for ri, row := range recs {
		ints, err := atoi(row)
		if err != nil {
			h.Error(
				w,
				fmt.Sprintf("Error: parsing CSV: record on line %d: %v", ri+1, err),
				h.StatusBadRequest)

			return
		}

		// Add the row to the transposed matrix as a column.
		for i, int := range ints {
			tran[i] = append(tran[i], int)
		}
	}

	// Build the response.
	var resp string
	for _, ints := range tran {
		resp += itos(ints) + "\n"
	}

	// The challenge description requires a trailing "\n" in the response.
	if len(resp) == 0 {
		resp += "\n"
	}

	fmt.Fprint(w, resp)
}
