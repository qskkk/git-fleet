package git

import (
	"context"
	"fmt"

	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/domain/repositories"
	"github.com/qskkk/git-fleet/internal/domain/services"
	"github.com/qskkk/git-fleet/internal/pkg/logger"
)

// StatusService implements the StatusService interface
type StatusService struct {
	gitRepo       repositories.GitRepository
	configService services.ConfigService
	logger        logger.Service
}

// NewStatusService creates a new status service
func NewStatusService(
	gitRepo repositories.GitRepository,
	configService services.ConfigService,
	logger logger.Service,
) services.StatusService {
	return &StatusService{
		gitRepo:       gitRepo,
		configService: configService,
		logger:        logger,
	}
}

// GetRepositoryStatus returns the status of a single repository
func (s *StatusService) GetRepositoryStatus(ctx context.Context, repoName string) (*entities.Repository, error) {
	s.logger.Debug(ctx, "Getting repository status", "repository", repoName)

	// Get the repository from config
	repo, err := s.configService.GetRepository(ctx, repoName)
	if err != nil {
		s.logger.Error(ctx, "Failed to get repository from config", err, "repository", repoName)
		return nil, fmt.Errorf("failed to get repository %s from config: %w", repoName, err)
	}

	updatedRepo, err := s.gitRepo.GetStatus(ctx, repo)
	if err != nil {
		s.logger.Error(ctx, "Failed to get repository status", err, "repository", repoName)
		return nil, fmt.Errorf("failed to get status for repository %s: %w", repoName, err)
	}

	return updatedRepo, nil
}

// GetGroupStatus returns the status of all repositories in a group
func (s *StatusService) GetGroupStatus(ctx context.Context, groupName string) ([]*entities.Repository, error) {
	s.logger.Info(ctx, "Getting group status", "group", groupName)

	// Get repositories for the group
	repos, err := s.configService.GetRepositoriesForGroups(ctx, []string{groupName})
	if err != nil {
		s.logger.Error(ctx, "Failed to get repositories for group", err, "group", groupName)
		return nil, fmt.Errorf("failed to get repositories for group %s: %w", groupName, err)
	}

	// Get status for each repository
	var statusRepos []*entities.Repository
	for _, repo := range repos {
		statusRepo, err := s.GetRepositoryStatus(ctx, repo.Name)
		if err != nil {
			// Continue with error information in the repository
			statusRepo = repo
			statusRepo.Status = "error"
		}
		statusRepos = append(statusRepos, statusRepo)
	}

	s.logger.Info(ctx, "Group status retrieved",
		"group", groupName,
		"repositories", len(statusRepos))

	return statusRepos, nil
}

// GetAllStatus returns the status of all repositories
func (s *StatusService) GetAllStatus(ctx context.Context) ([]*entities.Repository, error) {
	s.logger.Info(ctx, "Getting status for all repositories")

	// Get all repositories
	repos, err := s.configService.GetAllRepositories(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all repositories", err)
		return nil, fmt.Errorf("failed to get all repositories: %w", err)
	}

	// Get status for each repository
	var statusRepos []*entities.Repository
	for _, repo := range repos {
		statusRepo, err := s.GetRepositoryStatus(ctx, repo.Name)
		if err != nil {
			// Continue with error information in the repository
			statusRepo = repo
			statusRepo.Status = "error"
		}
		statusRepos = append(statusRepos, statusRepo)
	}

	s.logger.Info(ctx, "All repositories status retrieved",
		"repositories", len(statusRepos))

	return statusRepos, nil
}

// GetMultiGroupStatus returns the status of repositories in multiple groups
func (s *StatusService) GetMultiGroupStatus(ctx context.Context, groupNames []string) ([]*entities.Repository, error) {
	s.logger.Info(ctx, "Getting multi-group status", "groups", groupNames)

	// Get repositories for all groups
	repos, err := s.configService.GetRepositoriesForGroups(ctx, groupNames)
	if err != nil {
		s.logger.Error(ctx, "Failed to get repositories for groups", err, "groups", groupNames)
		return nil, fmt.Errorf("failed to get repositories for groups: %w", err)
	}

	// Get status for each repository
	var statusRepos []*entities.Repository
	for _, repo := range repos {
		statusRepo, err := s.GetRepositoryStatus(ctx, repo.Name)
		if err != nil {
			// Continue with error information in the repository
			statusRepo = repo
			statusRepo.Status = "error"
		}
		statusRepos = append(statusRepos, statusRepo)
	}

	s.logger.Info(ctx, "Multi-group status retrieved",
		"groups", groupNames,
		"repositories", len(statusRepos))

	return statusRepos, nil
}

// RefreshStatus refreshes the status of repositories
func (s *StatusService) RefreshStatus(ctx context.Context, repos []*entities.Repository) error {
	s.logger.Info(ctx, "Refreshing repository status", "repositories", len(repos))

	for _, repo := range repos {
		_, err := s.gitRepo.GetStatus(ctx, repo)
		if err != nil {
			s.logger.Error(ctx, "Failed to refresh repository status", err, "repository", repo.Name)
			// Continue with other repositories
		}
	}

	s.logger.Info(ctx, "Repository status refresh completed")
	return nil
}

// ValidateRepository validates if a repository is properly configured
func (s *StatusService) ValidateRepository(ctx context.Context, repo *entities.Repository) error {
	if repo == nil {
		return fmt.Errorf("repository cannot be nil")
	}

	if repo.Name == "" {
		return fmt.Errorf("repository name cannot be empty")
	}

	if repo.Path == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	// Check if repository path exists and is a valid Git repository
	if !s.gitRepo.IsValidDirectory(ctx, repo.Path) {
		return fmt.Errorf("repository path does not exist or is not accessible: %s", repo.Path)
	}

	if !s.gitRepo.IsValidRepository(ctx, repo.Path) {
		return fmt.Errorf("path is not a valid Git repository: %s", repo.Path)
	}

	return nil
}
