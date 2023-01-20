package main

import (
	"log"

	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
)

var config = conf.Config{}

func init() {
	options, err := parseOpts()
	if err != nil {
		log.Fatal(err)
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
	r := setupRouter()

	r.Run(config.HTTP.Host + ":" + config.HTTP.Port)
}
