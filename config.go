package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	loaded     bool
	path       string
	AdGuardBin string `yaml:"adGuardBin"`
}

func makeConfig(path string) *AppConfig {
	return &AppConfig{
		loaded:     false,
		path:       path,
		AdGuardBin: "adguardvpn-cli",
	}
}

func (c *AppConfig) load() error {
	data, err := os.ReadFile(c.path)

	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return err
	}

	err = yaml.Unmarshal(data, c)

	return err
}

func (c *AppConfig) save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.path, data, 0664)
	if err != nil {
		return err
	}
	return nil
}
