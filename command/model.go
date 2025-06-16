package command

import "github.com/qskkk/git-fleet/config"

var Handled = map[string]func(group string) (out string, err error){
	"status": ExecuteStatus,
	"ls":     ExecuteStatus,
	"pull":   ExecutePull,
	"pl":     ExecutePull,
}

var GlobalHandled = map[string]func(group string) (out string, err error){
	"config":    config.ExecuteConfig,
	"--config":  config.ExecuteConfig,
	"-c":        config.ExecuteConfig,
	"help":      ExecuteHelp,
	"--help":    ExecuteHelp,
	"-h":        ExecuteHelp,
	"version":   config.ExecuteVersionConfig,
	"--version": config.ExecuteVersionConfig,
	"-v":        config.ExecuteVersionConfig,
	"status":    ExecuteStatus,
	"--status":  ExecuteStatus,
	"-s":        ExecuteStatus,
}
