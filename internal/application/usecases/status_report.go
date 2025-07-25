package usecases

import (
	"context"
	"strings"

	"github.com/qskkk/git-fleet/v2/internal/application/ports/output"
	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/domain/services"
	"github.com/qskkk/git-fleet/v2/internal/pkg/errors"
)

// StatusReportUseCase handles repository status reporting
type StatusReportUseCase struct {
	configRepo    repositories.ConfigRepository
	gitRepo       repositories.GitRepository
	configService services.ConfigService
	statusService services.StatusService
	logger        services.LoggingService
	presenter     output.PresenterPort
}

// NewStatusReportUseCase creates a new StatusReportUseCase
func NewStatusReportUseCase(
	configRepo repositories.ConfigRepository,
	gitRepo repositories.GitRepository,
	configService services.ConfigService,
	statusService services.StatusService,
	logger services.LoggingService,
	presenter output.PresenterPort,
) *StatusReportUseCase {
	return &StatusReportUseCase{
		configRepo:    configRepo,
		gitRepo:       gitRepo,
		configService: configService,
		statusService: statusService,
		logger:        logger,
		presenter:     presenter,
	}
}

// StatusReportInput represents input for status reporting
type StatusReportInput struct {
	Groups      []string `json:"groups,omitempty"`
	Repository  string   `json:"repository,omitempty"`
	ShowAll     bool     `json:"show_all"`
	Refresh     bool     `json:"refresh"`
	ShowDetails bool     `json:"show_details"`
}

// StatusReportOutput represents output from status reporting
type StatusReportOutput struct {
	Repositories    []*entities.Repository `json:"repositories"`
	FormattedOutput string                 `json:"formatted_output"`
	Summary         *StatusSummary         `json:"summary"`
}

// StatusSummary represents a summary of repository statuses
type StatusSummary struct {
	TotalRepositories    int `json:"total_repositories"`
	CleanRepositories    int `json:"clean_repositories"`
	ModifiedRepositories int `json:"modified_repositories"`
	ErrorRepositories    int `json:"error_repositories"`
	WarningRepositories  int `json:"warning_repositories"`
}

// GetStatus gets the status of repositories
func (uc *StatusReportUseCase) GetStatus(ctx context.Context, input *StatusReportInput) (*StatusReportOutput, error) {
	uc.logger.Info(ctx, "Getting repository status", "input", input)

	var repositories []*entities.Repository
	var err error

	// Determine which repositories to check
	switch {
	case input.Repository != "":
		// Single repository
		repo, err := uc.statusService.GetRepositoryStatus(ctx, input.Repository)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get repository status", err, "repository", input.Repository)
			return nil, errors.WrapRepositoryNotFound(input.Repository)
		}
		repositories = []*entities.Repository{repo}

	case len(input.Groups) > 0:
		// Specific groups
		for _, groupName := range input.Groups {
			groupRepos, err := uc.statusService.GetGroupStatus(ctx, groupName)
			if err != nil {
				uc.logger.Error(ctx, "Failed to get group status", err, "group", groupName)
				return nil, errors.WrapGroupNotFound(groupName)
			}
			repositories = append(repositories, groupRepos...)
		}

	case input.ShowAll:
		// All repositories
		repositories, err = uc.statusService.GetAllStatus(ctx)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get all status", err)
			return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
		}

	default:
		// Default to all repositories
		repositories, err = uc.statusService.GetAllStatus(ctx)
		if err != nil {
			uc.logger.Error(ctx, "Failed to get all status", err)
			return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
		}
	}

	// Refresh status if requested
	if input.Refresh {
		uc.logger.Info(ctx, "Refreshing repository status")
		if err := uc.statusService.RefreshStatus(ctx, repositories); err != nil {
			uc.logger.Error(ctx, "Failed to refresh status", err)
			return nil, errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
		}
	}

	// Create summary
	summary := uc.createSummary(repositories)

	// Format output
	groupFilter := ""
	if len(input.Groups) > 0 {
		groupFilter = strings.Join(input.Groups, ", ")
	} else if input.Repository != "" {
		groupFilter = input.Repository
	}

	formattedOutput, err := uc.presenter.PresentStatus(ctx, repositories, groupFilter)
	if err != nil {
		uc.logger.Error(ctx, "Failed to format status output", err)
		// Don't fail the entire operation for formatting errors
		formattedOutput = "Error formatting status output"
	}

	uc.logger.Info(ctx, "Status report completed",
		"total", summary.TotalRepositories,
		"clean", summary.CleanRepositories,
		"modified", summary.ModifiedRepositories,
		"errors", summary.ErrorRepositories)

	return &StatusReportOutput{
		Repositories:    repositories,
		FormattedOutput: formattedOutput,
		Summary:         summary,
	}, nil
}

// createSummary creates a summary from repository statuses
func (uc *StatusReportUseCase) createSummary(repositories []*entities.Repository) *StatusSummary {
	summary := &StatusSummary{
		TotalRepositories: len(repositories),
	}

	for _, repo := range repositories {
		switch repo.Status {
		case entities.StatusClean:
			summary.CleanRepositories++
		case entities.StatusModified:
			summary.ModifiedRepositories++
		case entities.StatusError:
			summary.ErrorRepositories++
		case entities.StatusWarning:
			summary.WarningRepositories++
		}
	}

	return summary
}

// GetRepository gets a specific repository status
func (uc *StatusReportUseCase) GetRepository(ctx context.Context, name string) (*entities.Repository, error) {
	return uc.statusService.GetRepositoryStatus(ctx, name)
}

// GetAllRepositories gets all repository statuses
func (uc *StatusReportUseCase) GetAllRepositories(ctx context.Context) ([]*entities.Repository, error) {
	return uc.statusService.GetAllStatus(ctx)
}
