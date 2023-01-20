package data

import (
	"log"
	"time"

	"gitlab.com/yum2npm/yum2npm/pkg/config"
	"gitlab.com/yum2npm/yum2npm/pkg/yumrepodata"
)

// Repodata[<repo>]
type Repodata map[string]yumrepodata.PrimaryMetadata

// Modules[<repo>][<module>][<stream>]
type Modules map[string]map[string]map[string]yumrepodata.Module

type Update struct {
	Repodata Repodata
	Modules  Modules
}

func FetchPeriodically(interval time.Duration, repos []config.Repo, c chan Update) {
	for {
		go func() {
			r, m := fetch(repos)
			c <- Update{r, m}
		}()
		time.Sleep(interval)
	}
}

func fetch(repos []config.Repo) (Repodata, Modules) {
	r := Repodata{}
	m := Modules{}

	for _, repo := range repos {
		log.Printf("Refreshing %s", repo.Name)
		repomd, err := yumrepodata.GetRepoMetadata(repo.Url)
		if err != nil {
			log.Print(err)
			continue
		}

		primary, err := yumrepodata.GetPrimary(repo.Url, repomd)
		if err != nil {
			log.Print(err)
			continue
		}

		r[repo.Name] = primary

		modules, err := yumrepodata.GetModules(repo.Url, repomd)
		if err != nil {
			log.Print(err)
			continue
		}

		m[repo.Name] = modules
	}

	return r, m
}
