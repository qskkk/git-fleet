package main

import (
	"context"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/v2/internal/application/usecases"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/config"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/git"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/cli"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/tui"
	"github.com/qskkk/git-fleet/v2/internal/pkg/logger"
)

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Check for verbose/debug flag early
	verbose := false
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" || arg == "-d" || arg == "--debug" {
			verbose = true
			break
		}
	}

	// Initialize logger with appropriate level
	var loggerService logger.Service
	if verbose {
		loggerService = logger.NewWithLevel(logger.DEBUG)
	} else {
		loggerService = logger.NewWithLevel(logger.WARN)
	}

	// Initialize UI components
	stylesService := styles.NewService(styles.ThemeFleetName)
	presenter := cli.NewPresenter(stylesService)

	// Handle basic CLI commands without configuration
	basicHandler := cli.NewBasicHandler(stylesService)

	_ = basicHandler.Execute(ctx, os.Args)

	if basicHandler.HasHandled() {
		os.Exit(0)
	}

	// Initialize configuration
	configRepo := config.NewRepository()
	configService := config.NewService(configRepo, loggerService)
	validationService := config.NewValidationService()

	// Load configuration
	if err := configService.LoadConfig(ctx); err != nil {
		log.Errorf("Configuration Error: %v", err)
		os.Exit(1)
	}

	stylesService.SetTheme(styles.GetThemeFromString(configService.GetTheme(ctx)))

	// Initialize Git repository
	gitRepo := git.NewRepository()
	executorRepo := git.NewExecutor(stylesService)

	// Initialize services
	executionService := git.NewExecutionService(gitRepo, executorRepo, configService, loggerService)
	statusService := git.NewStatusService(gitRepo, configService, loggerService)

	// Initialize use cases
	executeCommandUC := usecases.NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		loggerService,
		presenter,
	)

	statusReportUC := usecases.NewStatusReportUseCase(
		configRepo,
		gitRepo,
		configService,
		statusService,
		loggerService,
		presenter,
	)

	manageConfigUC := usecases.NewManageConfigUseCase(
		configRepo,
		configService,
		validationService,
		loggerService,
		presenter,
	)

	// Determine if we should run in interactive mode
	if len(os.Args) == 1 {
		// Interactive mode
		runInteractiveMode(ctx, executeCommandUC, statusReportUC, manageConfigUC, stylesService, loggerService)
	} else {
		// CLI mode
		runCLIMode(ctx, os.Args, executeCommandUC, statusReportUC, manageConfigUC, stylesService, loggerService, verbose)
	}
}

// runInteractiveMode starts the interactive terminal UI
func runInteractiveMode(
	ctx context.Context,
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
	stylesService styles.Service,
	logger logger.Service,
) {
	logger.Info(ctx, "Starting interactive mode")

	// Create TUI
	tuiHandler := tui.NewHandler(executeCommandUC, statusReportUC, manageConfigUC, stylesService)

	// Run TUI
	if err := tuiHandler.Run(ctx); err != nil {
		log.Errorf("Terminal UI Error: %v", err)
		os.Exit(1)
	}
}

// runCLIMode handles command line interface
func runCLIMode(
	ctx context.Context,
	args []string,
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC *usecases.ManageConfigUseCase,
	stylesService styles.Service,
	logger logger.Service,
	verbose bool,
) {
	logLevel := "WARN"
	if verbose {
		logLevel = "DEBUG"
	}
	logger.Info(ctx, "Starting CLI mode", "args", args, "log_level", logLevel)

	// Create CLI handler
	cliHandler := cli.NewHandler(executeCommandUC, statusReportUC, manageConfigUC, stylesService)

	// Parse and execute command
	if err := cliHandler.Execute(ctx, args); err != nil {
		log.Errorf("Command Execution Error: %v", err)
		os.Exit(1)
	}
}
