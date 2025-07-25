package git

import (
	"context"
	"testing"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/domain/repositories"
	"github.com/qskkk/git-fleet/v2/internal/domain/services"
	"github.com/qskkk/git-fleet/v2/internal/pkg/errors"
	"github.com/qskkk/git-fleet/v2/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewStatusService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		gitRepo       repositories.GitRepository
		configService services.ConfigService
		logger        logger.Service
		wantNil       bool
	}{
		{
			name:          "should create status service with valid dependencies",
			gitRepo:       repositories.NewMockGitRepository(ctrl),
			configService: services.NewMockConfigService(ctrl),
			logger:        logger.NewMockService(ctrl),
			wantNil:       false,
		},
		{
			name:          "should create status service with nil git repository",
			gitRepo:       nil,
			configService: services.NewMockConfigService(ctrl),
			logger:        logger.NewMockService(ctrl),
			wantNil:       false,
		},
		{
			name:          "should create status service with nil config service",
			gitRepo:       repositories.NewMockGitRepository(ctrl),
			configService: nil,
			logger:        logger.NewMockService(ctrl),
			wantNil:       false,
		},
		{
			name:          "should create status service with nil logger",
			gitRepo:       repositories.NewMockGitRepository(ctrl),
			configService: services.NewMockConfigService(ctrl),
			logger:        nil,
			wantNil:       false,
		},
		{
			name:          "should create status service with all nil dependencies",
			gitRepo:       nil,
			configService: nil,
			logger:        nil,
			wantNil:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := NewStatusService(tt.gitRepo, tt.configService, tt.logger)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Implements(t, (*services.StatusService)(nil), result)

				// Verify the service is properly initialized
				statusService, ok := result.(*StatusService)
				assert.True(t, ok)
				assert.Equal(t, tt.gitRepo, statusService.gitRepo)
				assert.Equal(t, tt.configService, statusService.configService)
				assert.Equal(t, tt.logger, statusService.logger)
			}
		})
	}
}

func TestStatusService_GetRepositoryStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name     string
		repoName string
		setup    func()
		wantErr  bool
	}{
		{
			name:     "should get repository status successfully",
			repoName: "test-repo",
			setup: func() {
				repo := &entities.Repository{
					Name: "test-repo",
					Path: "/test/path",
				}
				updatedRepo := &entities.Repository{
					Name:   "test-repo",
					Path:   "/test/path",
					Status: entities.StatusClean,
					Branch: "main",
				}

				mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", "test-repo")
				mockConfigService.EXPECT().GetRepository(ctx, "test-repo").Return(repo, nil)
				mockGitRepo.EXPECT().GetStatus(ctx, repo).Return(updatedRepo, nil)
			},
			wantErr: false,
		},
		{
			name:     "should return error when repository not found in config",
			repoName: "unknown-repo",
			setup: func() {
				mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", "unknown-repo")
				mockConfigService.EXPECT().GetRepository(ctx, "unknown-repo").Return(nil, errors.ErrRepositoryNotFound)
				mockLogger.EXPECT().Error(ctx, "Failed to get repository from config", errors.ErrRepositoryNotFound, "repository", "unknown-repo")
			},
			wantErr: true,
		},
		{
			name:     "should return error when git status fails",
			repoName: "test-repo",
			setup: func() {
				repo := &entities.Repository{
					Name: "test-repo",
					Path: "/test/path",
				}

				mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", "test-repo")
				mockConfigService.EXPECT().GetRepository(ctx, "test-repo").Return(repo, nil)
				mockGitRepo.EXPECT().GetStatus(ctx, repo).Return(nil, errors.ErrGitStatusError)
				mockLogger.EXPECT().Error(ctx, "Failed to get repository status", gomock.Any(), "repository", "test-repo")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.GetRepositoryStatus(ctx, tt.repoName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestStatusService_GetGroupStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name      string
		groupName string
		setup     func()
		wantErr   bool
		wantLen   int
	}{
		{
			name:      "should get group status successfully",
			groupName: "test-group",
			setup: func() {
				repos := []*entities.Repository{
					{Name: "repo1", Path: "/path1"},
					{Name: "repo2", Path: "/path2"},
				}

				mockLogger.EXPECT().Info(ctx, "Getting group status", "group", "test-group")
				mockConfigService.EXPECT().GetRepositoriesForGroups(ctx, []string{"test-group"}).Return(repos, nil)

				// Mock GetRepository calls for each repo
				for _, repo := range repos {
					mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", repo.Name)
					mockConfigService.EXPECT().GetRepository(ctx, repo.Name).Return(repo, nil)
					updatedRepo := &entities.Repository{
						Name:   repo.Name,
						Path:   repo.Path,
						Status: entities.StatusClean,
					}
					mockGitRepo.EXPECT().GetStatus(ctx, repo).Return(updatedRepo, nil)
				}

				mockLogger.EXPECT().Info(ctx, "Group status retrieved", "group", "test-group", "repositories", 2)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:      "should return error when group not found",
			groupName: "unknown-group",
			setup: func() {
				mockLogger.EXPECT().Info(ctx, "Getting group status", "group", "unknown-group")
				mockConfigService.EXPECT().GetRepositoriesForGroups(ctx, []string{"unknown-group"}).Return(nil, errors.ErrGroupNotFound)
				mockLogger.EXPECT().Error(ctx, "Failed to get repositories for group", errors.ErrGroupNotFound, "group", "unknown-group")
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup()

			result, err := service.GetGroupStatus(ctx, tt.groupName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestStatusService_GetAllStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "should get all status successfully",
			setup: func() {
				repos := []*entities.Repository{
					{Name: "repo1", Path: "/path1"},
					{Name: "repo2", Path: "/path2"},
				}

				mockLogger.EXPECT().Info(ctx, "Getting status for all repositories")
				mockConfigService.EXPECT().GetAllRepositories(ctx).Return(repos, nil)

				// Mock GetRepository calls for each repo
				for _, repo := range repos {
					mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", repo.Name)
					mockConfigService.EXPECT().GetRepository(ctx, repo.Name).Return(repo, nil)
					updatedRepo := &entities.Repository{
						Name:   repo.Name,
						Path:   repo.Path,
						Status: entities.StatusClean,
					}
					mockGitRepo.EXPECT().GetStatus(ctx, repo).Return(updatedRepo, nil)
				}

				mockLogger.EXPECT().Info(ctx, "All repositories status retrieved", "repositories", 2)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "should return error when getting all repositories fails",
			setup: func() {
				mockLogger.EXPECT().Info(ctx, "Getting status for all repositories")
				mockConfigService.EXPECT().GetAllRepositories(ctx).Return(nil, errors.ErrRepositoryNotFound)
				mockLogger.EXPECT().Error(ctx, "Failed to get all repositories", errors.ErrRepositoryNotFound)
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			result, err := service.GetAllStatus(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestStatusService_ValidateRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    *entities.Repository
		setup   func()
		wantErr bool
	}{
		{
			name: "should validate repository successfully",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/valid/path",
			},
			setup: func() {
				mockGitRepo.EXPECT().IsValidDirectory(ctx, "/valid/path").Return(true)
				mockGitRepo.EXPECT().IsValidRepository(ctx, "/valid/path").Return(true)
			},
			wantErr: false,
		},
		{
			name:    "should return error for nil repository",
			repo:    nil,
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "should return error for empty repository name",
			repo: &entities.Repository{
				Name: "",
				Path: "/valid/path",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "should return error for empty repository path",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "should return error for invalid directory",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/invalid/path",
			},
			setup: func() {
				mockGitRepo.EXPECT().IsValidDirectory(ctx, "/invalid/path").Return(false)
			},
			wantErr: true,
		},
		{
			name: "should return error for invalid git repository",
			repo: &entities.Repository{
				Name: "test-repo",
				Path: "/valid/path",
			},
			setup: func() {
				mockGitRepo.EXPECT().IsValidDirectory(ctx, "/valid/path").Return(true)
				mockGitRepo.EXPECT().IsValidRepository(ctx, "/valid/path").Return(false)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup()

			err := service.ValidateRepository(ctx, tt.repo)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusService_RefreshStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		repos   []*entities.Repository
		setup   func()
		wantErr bool
	}{
		{
			name: "should refresh status successfully",
			repos: []*entities.Repository{
				{Name: "repo1", Path: "/path1"},
				{Name: "repo2", Path: "/path2"},
			},
			setup: func() {
				mockLogger.EXPECT().Info(ctx, "Refreshing repository status", "repositories", 2)
				mockGitRepo.EXPECT().GetStatus(ctx, gomock.Any()).Return(&entities.Repository{}, nil).Times(2)
				mockLogger.EXPECT().Info(ctx, "Repository status refresh completed")
			},
			wantErr: false,
		},
		{
			name: "should continue on error and complete refresh",
			repos: []*entities.Repository{
				{Name: "repo1", Path: "/path1"},
				{Name: "repo2", Path: "/path2"},
			},
			setup: func() {
				mockLogger.EXPECT().Info(ctx, "Refreshing repository status", "repositories", 2)
				mockGitRepo.EXPECT().GetStatus(ctx, gomock.Any()).Return(nil, errors.ErrGitStatusError)
				mockLogger.EXPECT().Error(ctx, "Failed to refresh repository status", errors.ErrGitStatusError, "repository", "repo1")
				mockGitRepo.EXPECT().GetStatus(ctx, gomock.Any()).Return(&entities.Repository{}, nil)
				mockLogger.EXPECT().Info(ctx, "Repository status refresh completed")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup()

			err := service.RefreshStatus(ctx, tt.repos)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusService_GetMultiGroupStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitRepo := repositories.NewMockGitRepository(ctrl)
	mockConfigService := services.NewMockConfigService(ctrl)
	mockLogger := logger.NewMockService(ctrl)

	service := NewStatusService(mockGitRepo, mockConfigService, mockLogger)
	ctx := context.Background()

	tests := []struct {
		name       string
		groupNames []string
		setup      func()
		wantErr    bool
		wantLen    int
	}{
		{
			name:       "should get multi-group status successfully",
			groupNames: []string{"group1", "group2"},
			setup: func() {
				repos := []*entities.Repository{
					{Name: "repo1", Path: "/path1"},
					{Name: "repo2", Path: "/path2"},
				}

				mockLogger.EXPECT().Info(ctx, "Getting multi-group status", "groups", []string{"group1", "group2"})
				mockConfigService.EXPECT().GetRepositoriesForGroups(ctx, []string{"group1", "group2"}).Return(repos, nil)

				// Mock GetRepository calls for each repo
				for _, repo := range repos {
					mockLogger.EXPECT().Debug(ctx, "Getting repository status", "repository", repo.Name)
					mockConfigService.EXPECT().GetRepository(ctx, repo.Name).Return(repo, nil)
					updatedRepo := &entities.Repository{
						Name:   repo.Name,
						Path:   repo.Path,
						Status: entities.StatusClean,
					}
					mockGitRepo.EXPECT().GetStatus(ctx, repo).Return(updatedRepo, nil)
				}

				mockLogger.EXPECT().Info(ctx, "Multi-group status retrieved", "groups", []string{"group1", "group2"}, "repositories", 2)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:       "should return error when groups not found",
			groupNames: []string{"unknown-group1", "unknown-group2"},
			setup: func() {
				mockLogger.EXPECT().Info(ctx, "Getting multi-group status", "groups", []string{"unknown-group1", "unknown-group2"})
				mockConfigService.EXPECT().GetRepositoriesForGroups(ctx, []string{"unknown-group1", "unknown-group2"}).Return(nil, errors.ErrGroupNotFound)
				mockLogger.EXPECT().Error(ctx, "Failed to get repositories for groups", errors.ErrGroupNotFound, "groups", []string{"unknown-group1", "unknown-group2"})
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup()

			result, err := service.GetMultiGroupStatus(ctx, tt.groupNames)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}
