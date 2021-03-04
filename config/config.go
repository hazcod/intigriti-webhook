package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	ConfigPath				string		`yaml:"-"`
	FindingIDs				[]string	`yaml:"findings"`

	IncludeNonReady       bool   `yaml:"include_non_ready"`
	CheckInterval         uint   `yaml:"check_interval_minutes"`

	Format				 string `yaml:"string"`

	WebhookURL            string `yaml:"webhook_url"`
	WebhookHeaders		  map[string]string `yaml:"webhook_headers"`

	IntigritiClientID     string `yaml:"intigriti_client_id"`
	IntigritiClientSecret string `yaml:"intigriti_client_secret"`
}

func ParseConfig(configPath string) (Config, error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "could not read configuration file")
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, errors.Wrap(err, "invalid yaml configuration")
	}

	return config, nil
}