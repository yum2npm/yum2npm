package main

import (
	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
	"gitlab.com/yum2npm/yum2npm/pkg/handlers"
	"gitlab.com/yum2npm/yum2npm/pkg/middleware"
	"net/http"
)

func setupRouter(config *conf.Config, repodata *data.Repodata, modules *data.Modules) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler(config.Repos))
	mux.HandleFunc("GET /repos", handlers.GetRepos(config.Repos))
	mux.HandleFunc("GET /repos/{repo}/packages", handlers.GetPackages(repodata))
	mux.HandleFunc("GET /repos/{repo}/packages/{package}", handlers.GetPackage(repodata))
	mux.HandleFunc("GET /repos/{repo}/modules", handlers.GetModules(modules))
	mux.HandleFunc("GET /repos/{repo}/modules/{module}/packages", handlers.GetModulePackages(modules))
	mux.HandleFunc("GET /repos/{repo}/modules/{module}/packages/:package", handlers.GetModulePackage(modules))

	return middleware.Log(mux)
}
