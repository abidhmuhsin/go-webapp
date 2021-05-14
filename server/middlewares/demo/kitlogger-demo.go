package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	customlogger "abidhmuhsin.com/gowebapp/server/middlewares"
	log "github.com/go-kit/kit/log"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!")
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", myHandler)

	var logger log.Logger
	// Logfmt is a structured, key=val logging format that is easy to read and parse
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// Direct any attempts to use Go's log package to our structured logger
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	// Log the timestamp (in UTC) and the callsite (file + line number) of the logging
	// call for debugging in the future.
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)

	// Create an instance of our LoggingMiddleware with our configured logger
	loggingMiddleware := customlogger.LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	// Start our HTTP server
	if err := http.ListenAndServe(":8000", loggedRouter); err != nil {
		logger.Log("status", "fatal", "err", err)
		os.Exit(1)
	}
}

/*
Run using go run kitlogger-demo.go
---Sample Output---

ts=2021-05-14T17:19:05.8726663Z loc=kitlogger-middleware.go:80 status=0 ip=127.0.0.1:63355 method=GET path=/ duration=0s
ts=2021-05-14T17:19:06.2761624Z loc=kitlogger-middleware.go:80 status=0 ip=127.0.0.1:63355 method=GET path=/favicon.ico duration=0s

*/
