package entities

import (
	"testing"
)

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name          string
		groupName     string
		repositories  []string
		expectedName  string
		expectedRepos []string
	}{
		{
			name:          "valid group creation",
			groupName:     "frontend",
			repositories:  []string{"webapp", "mobile"},
			expectedName:  "frontend",
			expectedRepos: []string{"webapp", "mobile"},
		},
		{
			name:          "group with single repository",
			groupName:     "backend",
			repositories:  []string{"api"},
			expectedName:  "backend",
			expectedRepos: []string{"api"},
		},
		{
			name:          "group with empty repositories",
			groupName:     "empty-group",
			repositories:  []string{},
			expectedName:  "empty-group",
			expectedRepos: []string{},
		},
		{
			name:          "group with nil repositories",
			groupName:     "nil-group",
			repositories:  nil,
			expectedName:  "nil-group",
			expectedRepos: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup(tt.groupName, tt.repositories)

			if group.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, group.Name)
			}

			if len(group.Repositories) != len(tt.expectedRepos) {
				t.Errorf("Expected %d repositories, got %d", len(tt.expectedRepos), len(group.Repositories))
			}

			for i, repo := range tt.expectedRepos {
				if i >= len(group.Repositories) || group.Repositories[i] != repo {
					t.Errorf("Expected repository %s at index %d, got %s", repo, i, group.Repositories[i])
				}
			}

			// Test default values
			if group.Description != "" {
				t.Errorf("Expected empty description by default, got %s", group.Description)
			}
		})
	}
}

func TestGroup_AddRepository(t *testing.T) {
	tests := []struct {
		name          string
		initialRepos  []string
		addRepo       string
		expectedRepos []string
	}{
		{
			name:          "add repository to empty group",
			initialRepos:  []string{},
			addRepo:       "new-repo",
			expectedRepos: []string{"new-repo"},
		},
		{
			name:          "add repository to existing group",
			initialRepos:  []string{"repo1", "repo2"},
			addRepo:       "repo3",
			expectedRepos: []string{"repo1", "repo2", "repo3"},
		},
		{
			name:          "add duplicate repository",
			initialRepos:  []string{"repo1", "repo2"},
			addRepo:       "repo1",
			expectedRepos: []string{"repo1", "repo2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup("test-group", tt.initialRepos)
			group.AddRepository(tt.addRepo)

			if len(group.Repositories) != len(tt.expectedRepos) {
				t.Errorf("Expected %d repositories, got %d", len(tt.expectedRepos), len(group.Repositories))
			}

			for i, expectedRepo := range tt.expectedRepos {
				if i >= len(group.Repositories) || group.Repositories[i] != expectedRepo {
					t.Errorf("Expected repository %s at index %d, got %s", expectedRepo, i, group.Repositories[i])
				}
			}
		})
	}
}

func TestGroup_RemoveRepository(t *testing.T) {
	tests := []struct {
		name          string
		initialRepos  []string
		removeRepo    string
		expectedRepos []string
	}{
		{
			name:          "remove repository from group",
			initialRepos:  []string{"repo1", "repo2", "repo3"},
			removeRepo:    "repo2",
			expectedRepos: []string{"repo1", "repo3"},
		},
		{
			name:          "remove first repository",
			initialRepos:  []string{"repo1", "repo2", "repo3"},
			removeRepo:    "repo1",
			expectedRepos: []string{"repo2", "repo3"},
		},
		{
			name:          "remove last repository",
			initialRepos:  []string{"repo1", "repo2", "repo3"},
			removeRepo:    "repo3",
			expectedRepos: []string{"repo1", "repo2"},
		},
		{
			name:          "remove non-existent repository",
			initialRepos:  []string{"repo1", "repo2"},
			removeRepo:    "repo3",
			expectedRepos: []string{"repo1", "repo2"},
		},
		{
			name:          "remove from empty group",
			initialRepos:  []string{},
			removeRepo:    "repo1",
			expectedRepos: []string{},
		},
		{
			name:          "remove only repository",
			initialRepos:  []string{"repo1"},
			removeRepo:    "repo1",
			expectedRepos: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup("test-group", tt.initialRepos)
			group.RemoveRepository(tt.removeRepo)

			if len(group.Repositories) != len(tt.expectedRepos) {
				t.Errorf("Expected %d repositories, got %d", len(tt.expectedRepos), len(group.Repositories))
			}

			for i, expectedRepo := range tt.expectedRepos {
				if i >= len(group.Repositories) || group.Repositories[i] != expectedRepo {
					t.Errorf("Expected repository %s at index %d, got %s", expectedRepo, i, group.Repositories[i])
				}
			}
		})
	}
}

func TestGroup_ContainsRepository(t *testing.T) {
	group := NewGroup("test-group", []string{"repo1", "repo2", "repo3"})

	tests := []struct {
		name     string
		repoName string
		expected bool
	}{
		{
			name:     "contains existing repository",
			repoName: "repo2",
			expected: true,
		},
		{
			name:     "contains first repository",
			repoName: "repo1",
			expected: true,
		},
		{
			name:     "contains last repository",
			repoName: "repo3",
			expected: true,
		},
		{
			name:     "does not contain non-existent repository",
			repoName: "repo4",
			expected: false,
		},
		{
			name:     "does not contain empty string",
			repoName: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := group.ContainsRepository(tt.repoName)
			if result != tt.expected {
				t.Errorf("Expected ContainsRepository(%s) to return %v, got %v", tt.repoName, tt.expected, result)
			}
		})
	}
}

func TestGroup_IsEmpty(t *testing.T) {
	tests := []struct {
		name         string
		repositories []string
		expected     bool
	}{
		{
			name:         "empty group",
			repositories: []string{},
			expected:     true,
		},
		{
			name:         "nil repositories",
			repositories: nil,
			expected:     true,
		},
		{
			name:         "group with one repository",
			repositories: []string{"repo1"},
			expected:     false,
		},
		{
			name:         "group with multiple repositories",
			repositories: []string{"repo1", "repo2", "repo3"},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup("test-group", tt.repositories)
			result := group.IsEmpty()
			if result != tt.expected {
				t.Errorf("Expected IsEmpty() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGroup_Count(t *testing.T) {
	tests := []struct {
		name         string
		repositories []string
		expected     int
	}{
		{
			name:         "empty group",
			repositories: []string{},
			expected:     0,
		},
		{
			name:         "nil repositories",
			repositories: nil,
			expected:     0,
		},
		{
			name:         "group with one repository",
			repositories: []string{"repo1"},
			expected:     1,
		},
		{
			name:         "group with multiple repositories",
			repositories: []string{"repo1", "repo2", "repo3"},
			expected:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup("test-group", tt.repositories)
			result := group.Count()
			if result != tt.expected {
				t.Errorf("Expected Count() to return %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGroup_Validate(t *testing.T) {
	tests := []struct {
		name         string
		group        *Group
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid group",
			group:       NewGroup("test-group", []string{"repo1", "repo2"}),
			expectError: false,
		},
		{
			name:         "group with empty name",
			group:        NewGroup("", []string{"repo1"}),
			expectError:  true,
			errorMessage: "group name cannot be empty",
		},
		{
			name:         "group with no repositories",
			group:        NewGroup("test-group", []string{}),
			expectError:  true,
			errorMessage: "group must contain at least one repository",
		},
		{
			name:         "group with nil repositories",
			group:        NewGroup("test-group", nil),
			expectError:  true,
			errorMessage: "group must contain at least one repository",
		},
		{
			name:         "group with empty name and no repositories",
			group:        NewGroup("", []string{}),
			expectError:  true,
			errorMessage: "group name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGroup_String(t *testing.T) {
	tests := []struct {
		name     string
		group    *Group
		expected string
	}{
		{
			name:     "group with repositories",
			group:    NewGroup("frontend", []string{"webapp", "mobile"}),
			expected: "Group{Name: frontend, Repositories: [webapp mobile]}",
		},
		{
			name:     "empty group",
			group:    NewGroup("empty", []string{}),
			expected: "Group{Name: empty, Repositories: []}",
		},
		{
			name:     "group with single repository",
			group:    NewGroup("backend", []string{"api"}),
			expected: "Group{Name: backend, Repositories: [api]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.group.String()
			if result != tt.expected {
				t.Errorf("Expected String() to return '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGroup_Description(t *testing.T) {
	group := NewGroup("test-group", []string{"repo1"})

	// Test default empty description
	if group.Description != "" {
		t.Errorf("Expected empty description by default, got '%s'", group.Description)
	}

	// Test setting description
	group.Description = "This is a test group"
	if group.Description != "This is a test group" {
		t.Errorf("Expected description 'This is a test group', got '%s'", group.Description)
	}
}

func TestGroup_ComplexOperations(t *testing.T) {
	// Test complex sequence of operations
	group := NewGroup("complex-group", []string{"repo1"})

	// Add multiple repositories
	group.AddRepository("repo2")
	group.AddRepository("repo3")
	group.AddRepository("repo2") // Duplicate - should not be added

	if group.Count() != 3 {
		t.Errorf("Expected 3 repositories after additions, got %d", group.Count())
	}

	// Remove middle repository
	group.RemoveRepository("repo2")

	if group.Count() != 2 {
		t.Errorf("Expected 2 repositories after removal, got %d", group.Count())
	}

	if group.ContainsRepository("repo2") {
		t.Error("Expected repo2 to be removed from group")
	}

	if !group.ContainsRepository("repo1") || !group.ContainsRepository("repo3") {
		t.Error("Expected repo1 and repo3 to still be in group")
	}

	// Validate should pass
	if err := group.Validate(); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}

	// Remove all repositories
	group.RemoveRepository("repo1")
	group.RemoveRepository("repo3")

	if !group.IsEmpty() {
		t.Error("Expected group to be empty after removing all repositories")
	}

	// Validation should fail now
	if err := group.Validate(); err == nil {
		t.Error("Expected validation to fail for empty group")
	}
}
