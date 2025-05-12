package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "pkg/config/config.yaml"

type Config struct {
	SiteTitle  string `yaml:"siteTitle"`
	ContentDir string `yaml:"contentDir"`
	OutputDir  string `yaml:"outputDir"`
	ServerPort string `yaml:"serverPort"`
}

var AppConfig Config

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{
		SiteTitle:  "Default Site Title",
		ContentDir: "content",
		OutputDir:  "output",
		ServerPort: "8080",
	}

	data, err := os.ReadFile(configPath)
	if err != nil {

		AppConfig = *cfg
		return cfg, nil
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	AppConfig = *cfg
	return cfg, nil
}

// SaveConfig saves the given configuration to the specified path.
func SaveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
