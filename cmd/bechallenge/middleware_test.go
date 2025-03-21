package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	// Prevent the app logging to the test log.
	l = slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestWebApiMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		payload    []byte
		wantBody   string
		wantStatus int
	}{
		{
			"huge-file",
			make([]byte, maxUploadSize+1),
			"Error: file upload size limit (10485760 bytes) exceeded\n",
			400,
		},
		{
			"empty-csv",
			[]byte{},
			"",
			200,
		},
		{
			"non-square-matrix",
			[]byte("1,2,3"),
			"Error: matrix is not square\n",
			400,
		},
		{
			"invalid-csv",
			[]byte("1,2,3\n4,3,5,7,8,9,"),
			"Error parsing CSV: record on line 2: wrong number of fields\n",
			400,
		},
	}

	h := webApiMiddleware(func(http.ResponseWriter, *http.Request) {})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}

func TestRecoverer(t *testing.T) {
	h := recoverer(func(http.ResponseWriter, *http.Request) {
		var zero int
		d := 10 / zero
		_ = d
	})

	r := httptest.NewRequest("POST", "/", nil)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != 500 {
		t.Errorf("Status code mismatch: got %d; want 500", w.Code)
	}
}
