package main

import (
	"bytes"
	"io"
	"log/slog"
	"mime/multipart"
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
		file       []byte
		wantBody   string
		wantStatus int
	}{
		{
			"file-too-large",
			make([]byte, maxUploadSize+1),
			"Error: file upload size limit (10485760 bytes) exceeded",
			400,
		},
		{
			"file-empty",
			[]byte{},
			"",
			200,
		},
		{
			"non-square-matrix",
			[]byte("1,2,3"),
			"Error: matrix is not square",
			400,
		},
	}

	h := webApiMiddleware(func(http.ResponseWriter, *http.Request) {})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)
			fWriter, err := writer.CreateFormFile("file", "file")
			if err != nil {
				t.Fatalf("unexpected multipart error %v", err)
			}

			_, err = fWriter.Write(tt.file)
			if err != nil {
				t.Fatalf("unexpected write error %v", err)
			}

			if err := writer.Close(); err != nil {
				t.Fatalf("unexpected writer close error %v", err)
			}

			r := httptest.NewRequest("POST", "/", buf)
			r.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Errorf(
					"Status code mismatch: got %d; want %d",
					w.Code, tt.wantStatus)
			}
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
