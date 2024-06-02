package main

import (
	"fmt"
	"log"
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
	mux := setupRouter()

	s := &http.Server{
		Addr:    config.HTTP.Host + ":" + config.HTTP.Port,
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
