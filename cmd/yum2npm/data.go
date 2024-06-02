package main

import "gitlab.com/yum2npm/yum2npm/pkg/data"

func receiveUpdates(c chan data.Update, repodata *data.Repodata, modules *data.Modules) {
	for u := range c {
		repodata = &u.Repodata
		modules = &u.Modules
	}
}
