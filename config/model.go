package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var configFile = os.ExpandEnv("$HOME/.config/git-fleet/.gfconfig.json")

type Repository struct {
	Path string `json:"path"`
}

type Config struct {
	Repositories map[string]Repository `json:"repositories"`
	Groups       map[string][]string   `json:"groups"`
}

var Cfg Config

func InitConfig() error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err = fmt.Errorf("configuration file not found. Please create it at $HOME/.config/.gfconfig.json.: %w", err)
		return err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		err := fmt.Errorf("it seems the configuration file is missing or unreadable. %w", err)
		return err
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		err := fmt.Errorf("error unmarshalling config file: %v\n", err)
		return err
	}

	return nil
}
