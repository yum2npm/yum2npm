package handlers

import (
	"bytes"
	_ "embed"
	"gitlab.com/yum2npm/yum2npm/pkg/utils"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
	"gitlab.com/yum2npm/yum2npm/pkg/yumrepodata"
)

type NpmIndex struct {
	Name     string                               `json:"name"`
	Versions map[string]yumrepodata.ModulePackage `json:"versions"`
}

//go:embed index.gohtml
var indexTemplate string

func IndexHandler(repos []config.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			utils.NotFound(w)
			return
		}
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		t, err := template.New("index.gohtml").Parse(indexTemplate)
		if err != nil {
			utils.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slog.Error("Error parsing index template", "Error", err)
			return
		}

		var out bytes.Buffer
		if err = t.Execute(&out, repos); err != nil {
			utils.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			slog.Error("Error executing index template", "Error", err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(out.Bytes())
		if err != nil {
			slog.Error("error while writing respone", "Error", err)
		}
	}
}

func GetRepos(repos []config.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.JsonResponse(w, repos)
	}
}

func GetPackages(repodata *data.Repodata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, exists := (*repodata)[r.PathValue("repo")]
		if !exists {
			utils.NotFound(w)
			return
		}

		utils.JsonResponse(w, repo)
	}
}

func GetModules(modules *data.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mods, exists := (*modules)[r.PathValue("repo")]
		if !exists {
			utils.NotFound(w)
			return
		}

		utils.JsonResponse(w, mods)
	}
}

func GetModulePackages(modules *data.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moduleIdentifier := strings.Split(r.PathValue("module"), ":")
		module, exists := (*modules)[r.PathValue("repo")][moduleIdentifier[0]][moduleIdentifier[1]]
		if !exists {
			utils.NotFound(w)
			return
		}

		utils.JsonResponse(w, module.GetPackages())
	}
}

func GetModulePackage(modules *data.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moduleIdentifier := strings.Split(r.PathValue("module"), ":")
		module, exists := (*modules)[r.PathValue("repo")][moduleIdentifier[0]][moduleIdentifier[1]]
		if !exists {
			utils.NotFound(w)
			return
		}

		packages, err := module.GetPackageVersions(r.PathValue("package"))
		if err != nil {
			utils.NotFound(w)
			return
		}

		index := NpmIndex{
			Name:     r.PathValue("package"),
			Versions: packages,
		}

		utils.JsonResponse(w, index)
	}
}

func GetPackage(repodata *data.Repodata) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, exists := (*repodata)[r.PathValue("repo")]
		if !exists {
			utils.NotFound(w)
			return
		}

		packages, err := repo.GetPackageVersions(r.PathValue("package"))
		if err != nil {
			utils.NotFound(w)
			return
		}

		index := NpmIndex{
			Name:     r.PathValue("package"),
			Versions: packages,
		}

		utils.JsonResponse(w, index)
	}
}
