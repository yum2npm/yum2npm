package main

import "gitlab.com/yum2npm/yum2npm/pkg/data"

var repodata = data.Repodata{}
var modules = data.Modules{}

func receiveUpdates(c chan data.Update) {
	for u := range c {
		repodata = u.Repodata
		modules = u.Modules
	}
}
