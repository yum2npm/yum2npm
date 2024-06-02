package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
)

var Version = "devel"

var config = conf.Config{}

func init() {
	options, err := parseOpts()
	if err != nil {
		log.Fatal(err)
	}

	if options.Version {
		fmt.Printf("yum2npm %s\n", Version)
		os.Exit(0)
	}

	config, err = conf.Init(options.Config)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan data.Update)
	go data.FetchPeriodically(config.RefreshInterval, config.Repos, c)
	go receiveUpdates(c)
}

func main() {
	mux := setupRouter(&config)

	server := &http.Server{
		Addr:    config.HTTP.ListenAddress,
		Handler: mux,
	}

	errCh := make(chan error, 1)

	go startServer(&config, server, errCh)
	go signalHandler(server, errCh)

	for {
		select {
		case err, ok := <-errCh:
			if !ok {
				return
			}
			if err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, context.DeadlineExceeded) {
				slog.Error("Error during server initialization", "Error", err)
				close(errCh)
			}
		}
	}
}
