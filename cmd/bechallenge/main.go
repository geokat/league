package main

import (
	"log/slog"
	h "net/http"
	"os"
)

type contextKey string

const (
	// in bytes
	maxUploadSize            = 10 * 1024 * 1024
	csvRecordsKey contextKey = "csvrecords"
)

var l *slog.Logger

func init() {
	// Include source code number in log entries.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)
	l = slog.Default()
}

func main() {
	mw := webApiMiddleware

	// Web API
	h.HandleFunc("/echo", mw(handleEcho))
	h.HandleFunc("/invert", mw(handleInvert))
	h.HandleFunc("/flatten", mw(handleFlatten))
	h.HandleFunc("/sum", mw(handleSum))
	h.HandleFunc("/multiply", mw(handleMultiply))

	// Stream API
	h.HandleFunc("/stream/echo", handleEchoStream)

	h.ListenAndServe(":8080", nil)
}
