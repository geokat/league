package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	h "net/http"
	"strings"
)

func handleEchoStream(w h.ResponseWriter, r *h.Request) {
	rdr := csv.NewReader(r.Body)
	for {
		row, err := rdr.Read()
		if err != nil {
			var pe *csv.ParseError

			// Not using http.Error() here as setting headers in
			// the middle of a stream is impossible.
			if err == io.EOF {
				break
			} else if errors.As(err, &pe) {
				fmt.Fprintln(w, "Error: parsing CSV: "+pe.Error())
			} else {
				l.Error("parsing CSV", "err", err)
				fmt.Fprintln(w, "Unexpected error")
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
