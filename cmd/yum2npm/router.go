package main

import (
	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
	"gitlab.com/yum2npm/yum2npm/pkg/handlers"
	"gitlab.com/yum2npm/yum2npm/pkg/middleware"
	"net/http"
	"net/http/pprof"
)

func setupRouter(config *conf.Config, profiling bool, repodata *data.Repodata, modules *data.Modules) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler(config.Repos))
	mux.HandleFunc("GET /repos", handlers.GetRepos(config.Repos))
	mux.HandleFunc("GET /repos/{repo}/packages", handlers.GetPackages(repodata))
	mux.HandleFunc("GET /repos/{repo}/packages/{package}", handlers.GetPackage(repodata))
	mux.HandleFunc("GET /repos/{repo}/modules", handlers.GetModules(modules))
	mux.HandleFunc("GET /repos/{repo}/modules/{module}/packages", handlers.GetModulePackages(modules))
	mux.HandleFunc("GET /repos/{repo}/modules/{module}/packages/:package", handlers.GetModulePackage(modules))

	if profiling {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	}

	return middleware.Log(mux)
}
