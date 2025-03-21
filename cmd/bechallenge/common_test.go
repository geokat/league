package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type formFileTestCase struct {
	name       string
	payload    []byte
	wantBody   string
	wantStatus int
}

// Helper; builds the test request for a form file upload, feeds it to
// the provided handler and asserts the response.
func runFormFileTestCase(
	t *testing.T,
	handler http.HandlerFunc,
	payload []byte,
	wantBody string,
	wantStatus int,
) {
	t.Helper()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	fWriter, err := writer.CreateFormFile("file", "file")
	if err != nil {
		t.Fatalf("unexpected multipart error %v", err)
	}

	if _, err = fWriter.Write(payload); err != nil {
		t.Fatalf("unexpected write error %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("unexpected writer close error %v", err)
	}

	r := httptest.NewRequest("POST", "/", buf)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	body := string(w.Body.Bytes())
	if body != wantBody {
		t.Errorf("Response body mismatch: got %q; want %q", body, wantBody)
	}
	if w.Code != wantStatus {
		t.Errorf("Status code mismatch: got %d; want %d", w.Code, wantStatus)
	}
}
