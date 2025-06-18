package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/qskkk/git-fleet/internal/domain/entities"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	if repo == nil {
		t.Error("NewRepository() should not return nil")
	}

	if _, ok := repo.(*Repository); !ok {
		t.Error("NewRepository() should return a *Repository")
	}
}

func TestRepository_GetStatus(t *testing.T) {
	repo := &Repository{}

	tests := []struct {
		name     string
		repoPath string
		wantErr  bool
	}{
		{
			name:     "non-existent path",
			repoPath: "/non/existent/path",
			wantErr:  false, // GetStatus handles errors gracefully
		},
		{
			name:     "empty path",
			repoPath: "",
			wantErr:  false, // GetStatus handles errors gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRepo := &entities.Repository{
				Name: "test-repo",
				Path: tt.repoPath,
			}

			ctx := context.Background()
			result, err := repo.GetStatus(ctx, testRepo)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result == nil {
				t.Error("GetStatus() should not return nil result")
			}

			// For invalid paths, IsValid should be false
			if tt.repoPath == "" || tt.repoPath == "/non/existent/path" {
				if result.IsValid {
					t.Error("GetStatus() should mark invalid paths as invalid")
				}
			}
		})
	}
}

func TestRepository_GetStatus_WithValidRepo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: tempDir,
	}

	ctx := context.Background()
	result, err := repo.GetStatus(ctx, testRepo)

	// For non-git directory, it should handle gracefully
	if err != nil {
		// This is expected for non-git directories
		if result == nil {
			t.Error("GetStatus() should return a result even on error")
		}
		return
	}

	if result == nil {
		t.Error("GetStatus() should not return nil result")
	}

	if result.Name != testRepo.Name {
		t.Errorf("GetStatus() result name = %v, want %v", result.Name, testRepo.Name)
	}
}

func TestRepository_GetLastCommit(t *testing.T) {
	repo := &Repository{}

	tests := []struct {
		name     string
		repoPath string
		wantErr  bool
	}{
		{
			name:     "non-existent path",
			repoPath: "/non/existent/path",
			wantErr:  true,
		},
		{
			name:     "empty path",
			repoPath: "",
			wantErr:  false, // May not error immediately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRepo := &entities.Repository{
				Name: "test-repo",
				Path: tt.repoPath,
			}

			ctx := context.Background()
			commit, err := repo.GetLastCommit(ctx, testRepo)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetLastCommit() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetLastCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if commit == nil {
				t.Error("GetLastCommit() should not return nil commit")
			}
		})
	}
}

func TestRepository_GetBranch(t *testing.T) {
	repo := &Repository{}

	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	ctx := context.Background()

	branch, err := repo.GetBranch(ctx, testRepo)

	// Should handle error gracefully
	if err == nil {
		t.Error("GetBranch() should return error for non-existent path")
	}

	if branch != "" {
		t.Error("GetBranch() should return empty string on error")
	}
}

func TestRepository_Fields(t *testing.T) {
	repo := &Repository{}

	// Test that we can create a Repository instance
	// This is a struct, so it can't be nil when created with &Repository{}
	_ = repo

	// Test that it has the expected methods available
	ctx := context.Background()
	testRepo := &entities.Repository{Name: "test", Path: "/tmp"}

	// Just ensure methods exist and can be called
	_, _ = repo.GetStatus(ctx, testRepo)
	_ = repo.IsValidDirectory(ctx, "/tmp")
}

func TestRepository_WithTempGitRepo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	gitDir := filepath.Join(tempDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: tempDir,
	}

	ctx := context.Background()

	// Test GetStatus with a git-like directory
	result, err := repo.GetStatus(ctx, testRepo)
	if result == nil {
		t.Error("GetStatus() should return a result")
	}

	// The error is expected since it's not a proper git repo
	if err != nil {
		t.Logf("Expected error for incomplete git repo: %v", err)
	}
}

func TestRepository_ContextCancellation(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/some/path",
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, _ := repo.GetStatus(ctx, testRepo)

	// GetStatus handles context gracefully but may not immediately error
	if result == nil {
		t.Error("GetStatus() should return a result even with cancelled context")
	}

	// The result should indicate an invalid repository
	if result != nil && result.IsValid {
		t.Error("GetStatus() should mark repository as invalid for cancelled context")
	}
}

func TestRepository_GetFileChanges(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	ctx := context.Background()
	created, modified, deleted, err := repo.GetFileChanges(ctx, testRepo)

	// Should handle error gracefully
	if err == nil {
		t.Error("GetFileChanges() should return error for non-existent path")
	}

	if created != 0 || modified != 0 || deleted != 0 {
		t.Error("GetFileChanges() should return zero values on error")
	}
}

func TestRepository_IsValidRepository(t *testing.T) {
	repo := &Repository{}
	ctx := context.Background()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "non-existent path",
			path: "/non/existent/path",
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.IsValidRepository(ctx, tt.path)
			if result != tt.want {
				t.Errorf("IsValidRepository() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestRepository_IsValidDirectory(t *testing.T) {
	repo := &Repository{}
	ctx := context.Background()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "non-existent path",
			path: "/non/existent/path",
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.IsValidDirectory(ctx, tt.path)
			if result != tt.want {
				t.Errorf("IsValidDirectory() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestRepository_ExecuteCommand(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	cmd := entities.NewGitCommand([]string{"status"})
	ctx := context.Background()

	result, err := repo.ExecuteCommand(ctx, testRepo, cmd)

	// Should return an error for non-existent repo, or handle it gracefully
	if err != nil && result == nil {
		t.Error("ExecuteCommand() should return a result even on error")
		return
	}

	if result != nil && result.Repository != testRepo.Name {
		t.Errorf("ExecuteCommand() result repository = %v, want %v", result.Repository, testRepo.Name)
	}
}

func TestRepository_HasUncommittedChanges(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	ctx := context.Background()
	hasChanges, err := repo.HasUncommittedChanges(ctx, testRepo)

	// Should handle error gracefully
	if err == nil {
		t.Error("HasUncommittedChanges() should return error for non-existent path")
	}

	if hasChanges {
		t.Error("HasUncommittedChanges() should return false on error")
	}
}

func TestRepository_GetAheadBehind(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	ctx := context.Background()
	ahead, behind, err := repo.GetAheadBehind(ctx, testRepo)

	// Should handle error gracefully - may or may not return error immediately
	if err != nil {
		// Error is expected, values should be zero
		if ahead != 0 || behind != 0 {
			t.Error("GetAheadBehind() should return zero values on error")
		}
	}
}

func TestRepository_GetRemotes(t *testing.T) {
	repo := &Repository{}
	testRepo := &entities.Repository{
		Name: "test-repo",
		Path: "/non/existent/path",
	}

	ctx := context.Background()
	remotes, err := repo.GetRemotes(ctx, testRepo)

	// Should handle error gracefully
	if err == nil {
		t.Error("GetRemotes() should return error for non-existent path")
	}

	if remotes != nil {
		t.Error("GetRemotes() should return nil remotes on error")
	}
}
