package main

import (
	"context"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/infrastructure/config"
	"github.com/qskkk/git-fleet/internal/infrastructure/git"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/cli"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/tui"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// LoggerAdapter adapts logger.Service to services.LoggingService
type LoggerAdapter struct {
	logger logger.Service
}

// Ensure LoggerAdapter implements services.LoggingService
var _ services.LoggingService = (*LoggerAdapter)(nil)

func (l *LoggerAdapter) Debug(ctx context.Context, message string, fields ...interface{}) {
	l.logger.Debug(ctx, message, fields...)
}

func (l *LoggerAdapter) Info(ctx context.Context, message string, fields ...interface{}) {
	l.logger.Info(ctx, message, fields...)
}

func (l *LoggerAdapter) Warn(ctx context.Context, message string, fields ...interface{}) {
	l.logger.Warn(ctx, message, fields...)
}

func (l *LoggerAdapter) Error(ctx context.Context, message string, err error, fields ...interface{}) {
	l.logger.Error(ctx, message, err, fields...)
}

func (l *LoggerAdapter) Fatal(ctx context.Context, message string, err error, fields ...interface{}) {
	l.logger.Fatal(ctx, message, err, fields...)
}

func (l *LoggerAdapter) SetLevel(level logger.Level) {
	l.logger.SetLevel(level)
}

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

	// Create logger adapter to match domain interface
	loggingService := &LoggerAdapter{logger: loggerService}

	// Initialize configuration
	configRepo := config.NewRepository()
	configService := config.NewService(configRepo, loggerService)
	validationService := config.NewValidationService()

	// Load configuration
	if err := configService.LoadConfig(ctx); err != nil {
		log.Errorf("Configuration Error: %v", err)
		os.Exit(1)
	}

	// Initialize Git repository
	gitRepo := git.NewRepository()
	executorRepo := git.NewExecutor()

	// Initialize services
	executionService := git.NewExecutionService(gitRepo, executorRepo, configService, loggingService)
	statusService := git.NewStatusService(gitRepo, configService, loggingService)

	// Initialize UI components
	stylesService := styles.NewService(configService.GetTheme())
	presenter := cli.NewPresenter(stylesService)

	// Initialize use cases
	executeCommandUC := usecases.NewExecuteCommandUseCase(
		configRepo,
		gitRepo,
		executorRepo,
		configService,
		executionService,
		validationService,
		loggingService,
		presenter,
	)

	statusReportUC := usecases.NewStatusReportUseCase(
		configRepo,
		gitRepo,
		configService,
		statusService,
		loggingService,
		presenter,
	)

	manageConfigUC := usecases.NewManageConfigUseCase(
		configRepo,
		configService,
		validationService,
		loggingService,
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
