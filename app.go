package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"abidhmuhsin.com/gowebapp/server"
)

func main() {

	handler := server.NewRouter()
	srv := &http.Server{
		Addr:         ":80",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server
	go func() {
		srv.ListenAndServe()
	}()
	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
