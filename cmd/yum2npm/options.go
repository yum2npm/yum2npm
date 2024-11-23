package main

import (
	flag "github.com/spf13/pflag"
	"os"
)

type Options struct {
	Config    string
	Version   bool
	Profiling bool
	Help      bool
}

func parseOpts() (Options, error) {
	var options Options

	flag.StringVarP(&options.Config, "config", "c", "/etc/yum2npm/config.yaml", "Path to config.yaml")
	flag.BoolVarP(&options.Version, "version", "v", false, "Print version information")
	flag.BoolVarP(&options.Profiling, "profiling", "p", false, "Enable profiling")
	flag.BoolVarP(&options.Help, "help", "h", false, "Show this help message")

	flag.Parse()

	if options.Help {
		flag.Usage()
		os.Exit(2)
	}

	return options, nil
}
