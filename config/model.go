package config

import (
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
