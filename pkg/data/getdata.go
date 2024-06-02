package data

import (
	"log/slog"
	"sync"
	"time"

	conf "gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/yumrepodata"
)

// Repodata [<repo>]
type Repodata map[string]yumrepodata.PrimaryMetadata

// Modules [<repo>][<module>][<stream>]
type Modules map[string]map[string]map[string]yumrepodata.Module

func FetchPeriodically(config *conf.Config, repodata *Repodata, modules *Modules) {
	for {
		mu := sync.Mutex{}
		mu.Lock()
		go func() {
			fetch(config, repodata, modules)
			mu.Unlock()
		}()

		time.Sleep(config.RefreshInterval)
	}
}

func fetch(config *conf.Config, repodata *Repodata, modules *Modules) {
	for _, repo := range config.Repos {
		slog.Info("Refreshing repository", "Repository", repo.Name)
		repomd, err := yumrepodata.GetRepoMetadata(repo.Url)
		if err != nil {
			slog.Error("Error fetching repository metadata", "Repository", repo.Name, "Error", err)
			continue
		}

		primary, err := yumrepodata.GetPrimary(repo.Url, repomd)
		if err != nil {
			slog.Error("Error fetching repository primary", "Repository", repo.Name, "Error", err)
			continue
		}

		(*repodata)[repo.Name] = primary

		(*modules)[repo.Name], err = yumrepodata.GetModules(repo.Url, repomd)
		if err != nil {
			slog.Error("Error fetching repository modules", "Repository", repo.Name, "Error", err)
			continue
		}
	}
}
