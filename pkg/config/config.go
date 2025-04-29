package config

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP            HTTP          `yaml:"http"`
	Repos           []Repo        `yaml:"repos"`
	RefreshInterval time.Duration `yaml:"refreshInterval"`
	Timeout         time.Duration `yaml:"timeout"`
}

type HTTP struct {
	ListenAddress string `yaml:"listen_addr"`
	CertFile      string `yaml:"cert_file"`
	KeyFile       string `yaml:"key_file"`
}

type Repo struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func Init(file string) (Config, error) {
	config := Config{
		HTTP: HTTP{
			ListenAddress: ":8080",
		},
		RefreshInterval: time.Hour,
		Timeout:         20 * time.Second,
	}

	configFile, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = configFile.Close() }()

	dec := yaml.NewDecoder(configFile)
	dec.KnownFields(true)

	if err := dec.Decode(&config); err != nil {
		return Config{}, err
	}

	if err := config.isValid(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) isValid() error {
	if len(c.HTTP.CertFile) > 0 && len(c.HTTP.KeyFile) == 0 {
		return fmt.Errorf("http key file must be specified if cert file is")
	}
	if len(c.HTTP.KeyFile) > 0 && len(c.HTTP.CertFile) == 0 {
		return fmt.Errorf("http cert file must be specified if key file is")
	}

	if len(c.Repos) == 0 {
		return fmt.Errorf("no repos are configured")
	}

	r := regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`)

	for i, repo := range c.Repos {
		if match := r.MatchString(repo.Name); !match {
			return fmt.Errorf("repo %d does not match regex \"^[a-zA-Z0-9-_]+$\" (%s)", i, repo.Name)
		}
		if len(repo.Url) == 0 {
			return fmt.Errorf("repo %d (%s) does not have an URL", i, repo.Name)
		}
		u, err := url.Parse(repo.Url)
		if err != nil {
			return fmt.Errorf("repo %d (%s) has an invalid URL: %w", i, repo.Name, err)
		}
		if u.Scheme == "" || u.Host == "" || u.Path == "" {
			return fmt.Errorf("repo %d (%s) has an invalid URL", i, repo.Name)
		}
	}

	return nil
}
