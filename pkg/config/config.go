package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	HTTP            HTTP          `yaml:"http"`
	Repos           []Repo        `yaml:"repos"`
	RefreshInterval time.Duration `yaml:"refreshInterval"`
}

type HTTP struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Repo struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func Init(file string) (Config, error) {
	config := Config{
		HTTP: HTTP{
			Host: "0.0.0.0",
			Port: "8080",
		},
		RefreshInterval: time.Hour,
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	if err := config.isValid(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) isValid() error {
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
	}

	return nil
}
