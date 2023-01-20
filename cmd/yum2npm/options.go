package main

import "github.com/jessevdk/go-flags"

type Options struct {
	Config string `long:"config" short:"c" default:"/etc/yum2npm/config.yaml" description:"Path to config.yaml"`
}

func parseOpts() (Options, error) {
	var options Options

	p := flags.NewParser(&options, flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		return options, err
	}

	return options, nil
}
