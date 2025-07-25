package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
)

func TestNewValidationService(t *testing.T) {
	service := NewValidationService()

	if service == nil {
		t.Fatal("NewValidationService() returned nil")
	}

	if _, ok := service.(*ValidationService); !ok {
		t.Error("NewValidationService() did not return a *ValidationService")
	}
}

func TestValidationService_ValidateRepository(t *testing.T) {
	ctx := context.Background()
	service := NewValidationService().(*ValidationService)

	t.Run("valid repository", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-repo",
			Path: "/absolute/path/to/repo",
		}

		err := service.ValidateRepository(ctx, repo)

		if err != nil {
			t.Errorf("ValidateRepository() error = %v, want nil", err)
		}
	})

	t.Run("nil repository", func(t *testing.T) {
		err := service.ValidateRepository(ctx, nil)

		if err == nil {
			t.Error("ValidateRepository() error = nil, want error")
		}

		expected := "repository cannot be nil"
		if err.Error() != expected {
			t.Errorf("ValidateRepository() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("empty repository name", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "",
			Path: "/absolute/path/to/repo",
		}

		err := service.ValidateRepository(ctx, repo)

		if err == nil {
			t.Error("ValidateRepository() error = nil, want error")
		}

		expected := "repository name cannot be empty"
		if err.Error() != expected {
			t.Errorf("ValidateRepository() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("empty repository path", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-repo",
			Path: "",
		}

		err := service.ValidateRepository(ctx, repo)

		if err == nil {
			t.Error("ValidateRepository() error = nil, want error")
		}

		expected := "repository path cannot be empty"
		if err.Error() != expected {
			t.Errorf("ValidateRepository() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("relative repository path", func(t *testing.T) {
		repo := &entities.Repository{
			Name: "test-repo",
			Path: "relative/path/to/repo",
		}

		err := service.ValidateRepository(ctx, repo)

		if err == nil {
			t.Error("ValidateRepository() error = nil, want error")
		}

		expectedPrefix := "path must be absolute"
		if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("ValidateRepository() error = %v, want prefix %v", err.Error(), expectedPrefix)
		}
	})
}

func TestValidationService_ValidateGroup(t *testing.T) {
	ctx := context.Background()
	service := NewValidationService().(*ValidationService)

	t.Run("valid group", func(t *testing.T) {
		group := entities.NewGroup("test-group", []string{"repo1", "repo2"})

		err := service.ValidateGroup(ctx, group)

		if err != nil {
			t.Errorf("ValidateGroup() error = %v, want nil", err)
		}
	})

	t.Run("nil group", func(t *testing.T) {
		err := service.ValidateGroup(ctx, nil)

		if err == nil {
			t.Error("ValidateGroup() error = nil, want error")
		}

		expected := "group cannot be nil"
		if err.Error() != expected {
			t.Errorf("ValidateGroup() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("group with empty name", func(t *testing.T) {
		group := entities.NewGroup("", []string{"repo1"})

		err := service.ValidateGroup(ctx, group)

		if err == nil {
			t.Error("ValidateGroup() error = nil, want error")
		}
	})

	t.Run("group with no repositories", func(t *testing.T) {
		group := entities.NewGroup("test-group", []string{})

		err := service.ValidateGroup(ctx, group)

		if err == nil {
			t.Error("ValidateGroup() error = nil, want error")
		}
	})
}

func TestValidationService_ValidateCommand(t *testing.T) {
	ctx := context.Background()
	service := NewValidationService().(*ValidationService)

	t.Run("valid command", func(t *testing.T) {
		cmd := entities.NewGitCommand([]string{"status"})

		err := service.ValidateCommand(ctx, cmd)

		if err != nil {
			t.Errorf("ValidateCommand() error = %v, want nil", err)
		}
	})

	t.Run("nil command", func(t *testing.T) {
		err := service.ValidateCommand(ctx, nil)

		if err == nil {
			t.Error("ValidateCommand() error = nil, want error")
		}

		expected := "command cannot be nil"
		if err.Error() != expected {
			t.Errorf("ValidateCommand() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("command with empty name", func(t *testing.T) {
		cmd := entities.NewGitCommand([]string{""})

		err := service.ValidateCommand(ctx, cmd)

		if err == nil {
			t.Error("ValidateCommand() error = nil, want error")
		}
	})

	t.Run("command with negative timeout", func(t *testing.T) {
		cmd := entities.NewGitCommand([]string{"status"})
		cmd.Timeout = -1

		err := service.ValidateCommand(ctx, cmd)

		if err == nil {
			t.Error("ValidateCommand() error = nil, want error")
		}
	})
}

func TestValidationService_ValidateConfig(t *testing.T) {
	ctx := context.Background()
	service := NewValidationService().(*ValidationService)

	t.Run("valid config", func(t *testing.T) {
		config := map[string]interface{}{
			"repositories": []string{"repo1", "repo2"},
		}

		err := service.ValidateConfig(ctx, config)

		if err != nil {
			t.Errorf("ValidateConfig() error = %v, want nil", err)
		}
	})

	t.Run("nil config", func(t *testing.T) {
		err := service.ValidateConfig(ctx, nil)

		if err == nil {
			t.Error("ValidateConfig() error = nil, want error")
		}

		expected := "configuration cannot be nil"
		if err.Error() != expected {
			t.Errorf("ValidateConfig() error = %v, want %v", err.Error(), expected)
		}
	})
}

func TestValidationService_ValidatePath(t *testing.T) {
	ctx := context.Background()
	service := NewValidationService().(*ValidationService)

	t.Run("empty path", func(t *testing.T) {
		err := service.ValidatePath(ctx, "")

		if err == nil {
			t.Error("ValidatePath() error = nil, want error")
		}

		expected := "path cannot be empty"
		if err.Error() != expected {
			t.Errorf("ValidatePath() error = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		err := service.ValidatePath(ctx, "relative/path")

		if err == nil {
			t.Error("ValidatePath() error = nil, want error")
		}

		expectedPrefix := "path must be absolute"
		if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("ValidatePath() error = %v, want prefix %v", err.Error(), expectedPrefix)
		}
	})

	t.Run("non-existent path", func(t *testing.T) {
		err := service.ValidatePath(ctx, "/non/existent/path")

		if err == nil {
			t.Error("ValidatePath() error = nil, want error")
		}

		expectedPrefix := "path does not exist"
		if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("ValidatePath() error = %v, want prefix %v", err.Error(), expectedPrefix)
		}
	})

	t.Run("valid existing directory", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "test_validation")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		err = service.ValidatePath(ctx, tempDir)

		if err != nil {
			t.Errorf("ValidatePath() error = %v, want nil", err)
		}
	})

	t.Run("path is a file not directory", func(t *testing.T) {
		// Create a temporary file for testing
		tempFile, err := os.CreateTemp("", "test_validation")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tempFile.Close()
		defer os.Remove(tempFile.Name())

		err = service.ValidatePath(ctx, tempFile.Name())

		if err == nil {
			t.Error("ValidatePath() error = nil, want error")
		}

		expectedPrefix := "path is not a directory"
		if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("ValidatePath() error = %v, want prefix %v", err.Error(), expectedPrefix)
		}
	})

	t.Run("path with permission issue", func(t *testing.T) {
		// This test is hard to simulate cross-platform
		// Skip it if we can't create a restricted directory
		tempDir, err := os.MkdirTemp("", "test_validation")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		restrictedPath := filepath.Join(tempDir, "restricted")
		err = os.Mkdir(restrictedPath, 0000) // No permissions
		if err != nil {
			t.Skip("Cannot create restricted directory for test")
		}

		// Try to access a subdirectory of the restricted directory
		testPath := filepath.Join(restrictedPath, "subdir")

		err = service.ValidatePath(ctx, testPath)

		if err == nil {
			t.Error("ValidatePath() error = nil, want error")
		}
	})
}
