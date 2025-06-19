package entities

import (
	"testing"
	"time"
)

func TestRepositoryStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   RepositoryStatus
		expected string
	}{
		{"Clean status", StatusClean, "Clean"},
		{"Modified status", StatusModified, "Modified"},
		{"Error status", StatusError, "Error"},
		{"Warning status", StatusWarning, "Warning"},
		{"Created status", StatusCreated, "Created"},
		{"Deleted status", StatusDeleted, "Deleted"},
		{"Unknown status", StatusUnknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.status))
			}
		})
	}
}

func TestRepository_HasChanges(t *testing.T) {
	tests := []struct {
		name     string
		repo     Repository
		expected bool
	}{
		{
			name:     "repository with created files",
			repo:     Repository{CreatedFiles: 1, ModifiedFiles: 0, DeletedFiles: 0},
			expected: true,
		},
		{
			name:     "repository with modified files",
			repo:     Repository{CreatedFiles: 0, ModifiedFiles: 1, DeletedFiles: 0},
			expected: true,
		},
		{
			name:     "repository with deleted files",
			repo:     Repository{CreatedFiles: 0, ModifiedFiles: 0, DeletedFiles: 1},
			expected: true,
		},
		{
			name:     "repository with multiple changes",
			repo:     Repository{CreatedFiles: 2, ModifiedFiles: 3, DeletedFiles: 1},
			expected: true,
		},
		{
			name:     "repository with no changes",
			repo:     Repository{CreatedFiles: 0, ModifiedFiles: 0, DeletedFiles: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.HasChanges()
			if result != tt.expected {
				t.Errorf("Expected HasChanges() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRepository_IsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		repo     Repository
		expected bool
	}{
		{
			name:     "healthy repository",
			repo:     Repository{IsValid: true, Status: StatusClean},
			expected: true,
		},
		{
			name:     "healthy modified repository",
			repo:     Repository{IsValid: true, Status: StatusModified},
			expected: true,
		},
		{
			name:     "invalid repository",
			repo:     Repository{IsValid: false, Status: StatusClean},
			expected: false,
		},
		{
			name:     "error repository",
			repo:     Repository{IsValid: true, Status: StatusError},
			expected: false,
		},
		{
			name:     "invalid and error repository",
			repo:     Repository{IsValid: false, Status: StatusError},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.IsHealthy()
			if result != tt.expected {
				t.Errorf("Expected IsHealthy() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRepository_GetDisplayPath(t *testing.T) {
	tests := []struct {
		name      string
		repo      Repository
		maxLength int
		expected  string
	}{
		{
			name:      "short path within limit",
			repo:      Repository{Path: "/short/path"},
			maxLength: 20,
			expected:  "/short/path",
		},
		{
			name:      "long path exceeding limit",
			repo:      Repository{Path: "/very/long/path/to/some/repository/folder"},
			maxLength: 20,
			expected:  "...repository/folder",
		},
		{
			name:      "path exactly at limit",
			repo:      Repository{Path: "/exact/length/path"},
			maxLength: 18,
			expected:  "/exact/length/path",
		},
		{
			name:      "empty path",
			repo:      Repository{Path: ""},
			maxLength: 10,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.GetDisplayPath(tt.maxLength)
			if result != tt.expected {
				t.Errorf("Expected GetDisplayPath() to return %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestRepository_UpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		repo           Repository
		expectedStatus RepositoryStatus
	}{
		{
			name: "invalid repository should be error",
			repo: Repository{
				IsValid:       false,
				CreatedFiles:  1,
				ModifiedFiles: 1,
				DeletedFiles:  1,
			},
			expectedStatus: StatusError,
		},
		{
			name: "valid repository with changes should be modified",
			repo: Repository{
				IsValid:       true,
				CreatedFiles:  1,
				ModifiedFiles: 0,
				DeletedFiles:  0,
			},
			expectedStatus: StatusModified,
		},
		{
			name: "valid repository with modified files should be modified",
			repo: Repository{
				IsValid:       true,
				CreatedFiles:  0,
				ModifiedFiles: 1,
				DeletedFiles:  0,
			},
			expectedStatus: StatusModified,
		},
		{
			name: "valid repository with deleted files should be modified",
			repo: Repository{
				IsValid:       true,
				CreatedFiles:  0,
				ModifiedFiles: 0,
				DeletedFiles:  1,
			},
			expectedStatus: StatusModified,
		},
		{
			name: "valid repository with no changes should be clean",
			repo: Repository{
				IsValid:       true,
				CreatedFiles:  0,
				ModifiedFiles: 0,
				DeletedFiles:  0,
			},
			expectedStatus: StatusClean,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.repo.UpdateStatus()
			if tt.repo.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, tt.repo.Status)
			}
		})
	}
}

func TestRepository_Fields(t *testing.T) {
	now := time.Now()
	repo := Repository{
		Name:          "test-repo",
		Path:          "/path/to/repo",
		Status:        StatusModified,
		Branch:        "feature-branch",
		CreatedFiles:  1,
		ModifiedFiles: 2,
		DeletedFiles:  3,
		LastChecked:   now,
		IsValid:       true,
		ErrorMessage:  "test error",
	}

	// Test all fields are set correctly
	if repo.Name != "test-repo" {
		t.Errorf("Expected Name 'test-repo', got %s", repo.Name)
	}

	if repo.Path != "/path/to/repo" {
		t.Errorf("Expected Path '/path/to/repo', got %s", repo.Path)
	}

	if repo.Status != StatusModified {
		t.Errorf("Expected Status %s, got %s", StatusModified, repo.Status)
	}

	if repo.Branch != "feature-branch" {
		t.Errorf("Expected Branch 'feature-branch', got %s", repo.Branch)
	}

	if repo.CreatedFiles != 1 {
		t.Errorf("Expected CreatedFiles 1, got %d", repo.CreatedFiles)
	}

	if repo.ModifiedFiles != 2 {
		t.Errorf("Expected ModifiedFiles 2, got %d", repo.ModifiedFiles)
	}

	if repo.DeletedFiles != 3 {
		t.Errorf("Expected DeletedFiles 3, got %d", repo.DeletedFiles)
	}

	if !repo.LastChecked.Equal(now) {
		t.Errorf("Expected LastChecked %v, got %v", now, repo.LastChecked)
	}

	if !repo.IsValid {
		t.Error("Expected IsValid to be true")
	}

	if repo.ErrorMessage != "test error" {
		t.Errorf("Expected ErrorMessage 'test error', got %s", repo.ErrorMessage)
	}
}

func TestRepository_RepositoryCreation(t *testing.T) {
	repo := Repository{
		Name: "test-repo",
		Path: "/path/to/repo",
	}

	// Test default values
	if repo.Status != "" {
		t.Errorf("Expected empty status by default, got %s", repo.Status)
	}

	if repo.Branch != "" {
		t.Errorf("Expected empty branch by default, got %s", repo.Branch)
	}

	if repo.CreatedFiles != 0 {
		t.Errorf("Expected 0 created files by default, got %d", repo.CreatedFiles)
	}

	if repo.ModifiedFiles != 0 {
		t.Errorf("Expected 0 modified files by default, got %d", repo.ModifiedFiles)
	}

	if repo.DeletedFiles != 0 {
		t.Errorf("Expected 0 deleted files by default, got %d", repo.DeletedFiles)
	}

	if repo.IsValid {
		t.Error("Expected IsValid to be false by default")
	}

	if repo.ErrorMessage != "" {
		t.Errorf("Expected empty ErrorMessage by default, got %s", repo.ErrorMessage)
	}
}
