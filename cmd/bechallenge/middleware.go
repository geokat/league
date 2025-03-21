package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	h "net/http"
	"runtime/debug"
)

// Does the prep work common to the handlers in our web API:
//   - prevent reading in unreasonable amounts of data
//   - make the payload available as CSV records
//   - handle panics in handler goroutines
//
// Reports known error messages to the user. Unexpected error messages
// only go in the logs as they can contain sensitive info about our infra.
func webApiMiddleware(next h.HandlerFunc) h.HandlerFunc {
	handler := func(w h.ResponseWriter, r *h.Request) {
		r.Body = h.MaxBytesReader(w, r.Body, maxUploadSize)

		f, _, err := r.FormFile("file")
		if err != nil {
			var mbe *h.MaxBytesError
			if errors.As(err, &mbe) {
				m := fmt.Sprintf(
					"Error: file upload size limit (%d bytes) exceeded",
					maxUploadSize)
				h.Error(w, m, h.StatusBadRequest)
			} else {
				l.Error("getting form file", "err", err)
				h.Error(w, "Error: unexpected error", h.StatusInternalServerError)
			}

			return
		}
		defer f.Close()

		recs, err := csv.NewReader(f).ReadAll()
		if err != nil {
			var pe *csv.ParseError
			if errors.As(err, &pe) {
				h.Error(w, "Error parsing CSV: "+pe.Error(), h.StatusBadRequest)
			} else {
				l.Error("parsing CSV", "err", err)
				h.Error(w, "Error: unexpected error", h.StatusInternalServerError)
			}

			return
		}

		// Make sure it's a square matrix.
		if len(recs) > 0 && len(recs) != len(recs[0]) {
			h.Error(w, "Error: matrix is not square", h.StatusBadRequest)

			return
		}

		// Make records available to downstream handlers.
		ctx := context.WithValue(r.Context(), csvRecordsKey, recs)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return recoverer(handler)
}

// Handles panics by logging the call trace and returning an error
// response to the user.
func recoverer(next h.HandlerFunc) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Error("Recovered from panic", "err", err, "trace", debug.Stack())
				h.Error(w, "Error: unexpected error", h.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}
}
