package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qskkk/git-fleet/internal/application/ports/output"
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// ExecuteCommandUseCase handles command execution business logic
type ExecuteCommandUseCase struct {
	configRepo        repositories.ConfigRepository
	gitRepo           repositories.GitRepository
	executorRepo      repositories.ExecutorRepository
	configService     services.ConfigService
	executionService  services.ExecutionService
	validationService services.ValidationService
	logger            services.LoggingService
	presenter         output.PresenterPort
}

// NewExecuteCommandUseCase creates a new ExecuteCommandUseCase
func NewExecuteCommandUseCase(
	configRepo repositories.ConfigRepository,
	gitRepo repositories.GitRepository,
	executorRepo repositories.ExecutorRepository,
	configService services.ConfigService,
	executionService services.ExecutionService,
	validationService services.ValidationService,
	logger services.LoggingService,
	presenter output.PresenterPort,
) *ExecuteCommandUseCase {
	return &ExecuteCommandUseCase{
		configRepo:        configRepo,
		gitRepo:           gitRepo,
		executorRepo:      executorRepo,
		configService:     configService,
		executionService:  executionService,
		validationService: validationService,
		logger:            logger,
		presenter:         presenter,
	}
}

// ExecuteCommandInput represents input for command execution
type ExecuteCommandInput struct {
	Groups       []string `json:"groups"`
	CommandStr   string   `json:"command"`
	Parallel     bool     `json:"parallel"`
	AllowFailure bool     `json:"allow_failure"`
	Timeout      int      `json:"timeout,omitempty"`
}

// ExecuteCommandOutput represents output from command execution
type ExecuteCommandOutput struct {
	Summary         *entities.Summary `json:"summary"`
	FormattedOutput string            `json:"formatted_output"`
	Success         bool              `json:"success"`
}

// Execute executes a command on specified groups
func (uc *ExecuteCommandUseCase) Execute(ctx context.Context, input *ExecuteCommandInput) (*ExecuteCommandOutput, error) {
	uc.logger.Info(ctx, "Starting command execution", "groups", input.Groups, "command", input.CommandStr)

	// Validate input
	if err := uc.validateInput(input); err != nil {
		uc.logger.Error(ctx, "Invalid input", err, "input", input)
		return nil, errors.WrapInvalidInput(err)
	}

	// Parse command
	command, err := uc.executionService.ParseCommand(ctx, input.CommandStr)
	if err != nil {
		uc.logger.Error(ctx, "Failed to parse command", err, "command", input.CommandStr)
		return nil, errors.WrapCommandParsingError(err)
	}

	// Apply timeout if specified
	if input.Timeout > 0 {
		command.Timeout = time.Duration(input.Timeout) * time.Second
	}
	command.AllowFailure = input.AllowFailure

	// Validate command
	if err := uc.validationService.ValidateCommand(ctx, command); err != nil {
		uc.logger.Error(ctx, "Invalid command", err, "command", command)
		return nil, errors.WrapInvalidInput(err)
	}

	// Check if it's a built-in command
	if uc.executionService.IsBuiltInCommand(command.Args[0]) {
		return uc.executeBuiltInCommand(ctx, command.Args[0], input.Groups)
	}

	// Get repositories for groups
	repositories, err := uc.configService.GetRepositoriesForGroups(ctx, input.Groups)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get repositories for groups", err, "groups", input.Groups)
		return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
	}

	if len(repositories) == 0 {
		uc.logger.Warn(ctx, "No repositories found for specified groups", "groups", input.Groups)
		return &ExecuteCommandOutput{
			Summary:         entities.NewSummary(),
			FormattedOutput: "No repositories found for specified groups",
			Success:         true,
		}, nil
	}

	// Execute command
	var summary *entities.Summary
	if input.Parallel {
		summary, err = uc.executorRepo.ExecuteInParallel(ctx, repositories, command)
	} else {
		summary, err = uc.executorRepo.ExecuteSequential(ctx, repositories, command)
	}

	if err != nil {
		uc.logger.Error(ctx, "Failed to execute command", err, "command", command, "repositories", len(repositories))
		return nil, errors.WrapFailedToExecuteCommand(err)
	}

	// Format output
	formattedOutput, err := uc.presenter.PresentSummary(ctx, summary)
	if err != nil {
		uc.logger.Error(ctx, "Failed to format output", err)
		// Don't fail the entire operation for formatting errors
		formattedOutput = "Error formatting output"
	}

	success := !summary.HasFailures()
	uc.logger.Info(
		ctx,
		fmt.Sprintf(
			"%s command execution completed \nSuccess: %t \nTotal Repositories: %d \nSuccessful Executions: %d \nFailed Executions: %d",
			command.GetFullCommand(),
			success,
			summary.TotalRepositories,
			summary.SuccessfulExecutions,
			summary.FailedExecutions,
		),
	)

	if uc.logger.GetLevel() <= logger.INFO {
		uc.logger.Debug(ctx, "Command execution summary", "summary", summary)
	}

	for _, result := range summary.Results {
		uc.logger.Debug(ctx, result.Output)
	}

	return &ExecuteCommandOutput{
		Summary:         summary,
		FormattedOutput: formattedOutput,
		Success:         success,
	}, nil
}

// executeBuiltInCommand handles built-in commands
func (uc *ExecuteCommandUseCase) executeBuiltInCommand(ctx context.Context, cmdName string, groups []string) (*ExecuteCommandOutput, error) {
	output, err := uc.executionService.ExecuteBuiltInCommand(ctx, cmdName, groups)
	if err != nil {
		return nil, err
	}

	summary := entities.NewSummary()
	summary.SuccessfulExecutions = 1
	summary.TotalRepositories = 1
	summary.Finalize()

	return &ExecuteCommandOutput{
		Summary:         summary,
		FormattedOutput: output,
		Success:         true,
	}, nil
}

// validateInput validates the command execution input
func (uc *ExecuteCommandUseCase) validateInput(input *ExecuteCommandInput) error {
	if len(input.Groups) == 0 {
		return errors.ErrAtLeastOneGroupRequired
	}

	if strings.TrimSpace(input.CommandStr) == "" {
		return errors.ErrCommandStringEmpty
	}

	if input.Timeout < 0 {
		return errors.ErrTimeoutCannotBeNegative
	}

	return nil
}

// GetAvailableCommands returns available commands
func (uc *ExecuteCommandUseCase) GetAvailableCommands(ctx context.Context) ([]string, error) {
	return uc.executionService.GetAvailableCommands(ctx)
}
