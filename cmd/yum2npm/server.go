package main

import (
	"context"
	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer(config *conf.Config, server *http.Server, errCh chan<- error) {
	if len(config.HTTP.CertFile) > 0 && len(config.HTTP.KeyFile) > 0 {
		errCh <- server.ListenAndServeTLS(config.HTTP.CertFile, config.HTTP.KeyFile)
	} else {
		errCh <- server.ListenAndServe()
	}
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
