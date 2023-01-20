package handlers

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/data"
	"gitlab.com/yum2npm/yum2npm/pkg/yumrepodata"
)

type NpmIndex struct {
	Name     string                               `json:"name"`
	Versions map[string]yumrepodata.ModulePackage `json:"versions"`
}

//go:embed index.html
var indexTemplate string

func IndexHandler(repos []config.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := template.New("index.html").Parse(indexTemplate)
		if err != nil {
			c.String(http.StatusNotFound, http.StatusText(500), 500)
			log.Print(err)
			return
		}

		var out bytes.Buffer
		if err = t.Execute(&out, repos); err != nil {
			c.String(http.StatusNotFound, http.StatusText(500), 500)
			log.Print(err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", out.Bytes())
	}
}

func GetRepos(repos []config.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, repos)
	}
}

func GetPackages(repodata *data.Repodata) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := (*repodata)[c.Param("repo")]; !exists {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}
		c.JSON(http.StatusOK, (*repodata)[c.Param("repo")])
	}
}

func GetModules(modules *data.Modules) gin.HandlerFunc {
	return func(c *gin.Context) {
		mods, exists := (*modules)[c.Param("repo")]
		if !exists {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}
		c.JSON(http.StatusOK, mods)
	}
}

func GetModulePackages(modules *data.Modules) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleIdentifier := strings.Split(c.Param("module"), ":")
		module, exists := (*modules)[c.Param("repo")][moduleIdentifier[0]][moduleIdentifier[1]]
		if !exists {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}

		c.JSON(http.StatusOK, module.GetPackages())
	}
}

func GetModulePackage(modules *data.Modules) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleIdentifier := strings.Split(c.Param("module"), ":")
		module, exists := (*modules)[c.Param("repo")][moduleIdentifier[0]][moduleIdentifier[1]]
		if !exists {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}

		packages, err := module.GetPackageVersions(c.Param("package"))
		if err != nil {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}

		index := NpmIndex{
			Name:     c.Param("package"),
			Versions: packages,
		}

		c.JSON(http.StatusOK, index)
	}
}

func GetPackage(repodata *data.Repodata) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo, exists := (*repodata)[c.Param("repo")]
		if !exists {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}

		packages, err := repo.GetPackageVersions(c.Param("package"))
		if err != nil {
			c.String(http.StatusNotFound, http.StatusText(404))
			return
		}

		index := NpmIndex{
			Name:     c.Param("package"),
			Versions: packages,
		}

		c.JSON(http.StatusOK, index)
	}
}
