package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	h "net/http"
	"strings"
)

// An example of a zero-allocation echo API handler where the payload
// is processed in stream while being uploaded. This way we can handle
// very large uploads without exhausting memory.
//
// Send request with:
//
//	curl -s -T matrix.csv "localhost:8080/stream/echo"
//
// Flatten, Sum and Multiply can also be implemented as stream APIs.
// Invert is not a good fit because it requires reading in all of the
// body before it can start producing output.
func handleEchoStream(w h.ResponseWriter, r *h.Request) {
	rdr := csv.NewReader(r.Body)
	for {
		row, err := rdr.Read()
		if err != nil {
			var pe *csv.ParseError

			if err == io.EOF {
				break
			} else if errors.As(err, &pe) {
				// Not using http.Error() as setting headers in the
				// middle of a stream is impossible.
				_, err := fmt.Fprintln(w, "Error: parsing CSV: "+pe.Error())
				if err != nil {
					l.Error("writing response error message", "error", err)
				}
			} else {
				l.Error("parsing CSV", "err", err)
				_, err := fmt.Fprintln(w, "Unexpected error")
				if err != nil {
					l.Error("writing response error message", "error", err)
				}
			}
			return
		}

		_, err = w.Write([]byte(strings.Join(row, ",") + "\n"))
		if err != nil {
			l.Error("writing response", "error", err)
			return
		}
	}
}
