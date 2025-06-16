package command

import (
	"bytes"
	"fmt"
	"time"

	"github.com/qskkk/git-fleet/config"
	"github.com/qskkk/git-fleet/style"
)

var Handled = map[string]func(group string) (out string, err error){
	"status": ExecuteStatus,
	"ls":     ExecuteStatus,
	"pull":   ExecutePull,
	"pl":     ExecutePull,
	"fetch":  ExecuteFetchAll,
	"fa":     ExecuteFetchAll,
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

type SummaryData struct {
	SuccessCount  int
	ErrorCount    int
	TargetGroup   string
	Command       string
	ExecutionTime time.Duration
}

func (sd *SummaryData) String() string {
	var result bytes.Buffer

	// Create summary table data
	summaryData := [][]string{
		{"✅ Successful Repositories", fmt.Sprintf("%d", sd.SuccessCount)},
		{"❌ Failed Repositories", fmt.Sprintf("%d", sd.ErrorCount)},
		{"🎯 Target Group", sd.TargetGroup},
		{"🔧 Command Executed", sd.Command},
		{"⌛ Execution Time", sd.ExecutionTime.String()},
	}

	result.WriteString(style.TitleStyle.Render("📊 Execution Summary") + "\n\n")
	summaryTable := style.CreateSummaryTable(summaryData)
	result.WriteString(summaryTable.String())

	return result.String()
}
