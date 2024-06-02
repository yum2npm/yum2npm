package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer(server *http.Server, errCh chan<- error) {
	errCh <- server.ListenAndServe()
}

func shutdownServer(server *http.Server, errCh chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer close(errCh)
	errCh <- server.Shutdown(ctx)
}

func signalHandler(server *http.Server, errCh chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	slog.Info("Shutting down server...")
	go shutdownServer(server, errCh)

	<-c
	slog.Info("Killing server...")
	os.Exit(1)
}
