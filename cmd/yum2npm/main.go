package main

import (
	"context"
	"errors"
	"fmt"
	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var Version = "devel"

func main() {
	options, err := parseOpts()
	if err != nil {
		log.Fatal(err)
	}

	if options.Version {
		fmt.Printf("yum2npm %s\n", Version)
		os.Exit(0)
	}

	config, err := conf.Init(options.Config)
	if err != nil {
		log.Fatal(err)
	}

	repodata := make(data.Repodata)
	modules := make(data.Modules)
	go data.FetchPeriodically(&config, &repodata, &modules)

	mux := setupRouter(&config, options.Profiling, &repodata, &modules)

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
