package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Version will be set at build time using ldflags
var Version = "dev"

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
		err = fmt.Errorf("❌ Configuration file not found.\n📁 Please create it at: %s", configFile)
		return err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		err := fmt.Errorf("❌ Configuration file is missing or unreadable: %w", err)
		return err
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		err := fmt.Errorf("❌ Invalid JSON in configuration file: %v", err)
		return err
	}

	return nil
}
