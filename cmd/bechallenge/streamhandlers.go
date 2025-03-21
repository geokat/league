package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	h "net/http"
	"strings"
)

// Not using http.Error() when handling errors here as setting headers
// in the middle of a stream is impossible.
func handleEchoStream(w h.ResponseWriter, r *h.Request) {
	rdr := csv.NewReader(r.Body)
	for {
		row, err := rdr.Read()
		if err != nil {
			var pe *csv.ParseError

			if err == io.EOF {
				break
			} else if errors.As(err, &pe) {
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
