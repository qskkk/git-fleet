package repositories

import (
	"testing"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
)

func TestConfig_GetRepository(t *testing.T) {
	config := &Config{
		Repositories: map[string]*RepositoryConfig{
			"repo1": {Path: "/path/to/repo1"},
			"repo2": {Path: "/path/to/repo2"},
		},
	}

	t.Run("existing repository", func(t *testing.T) {
		repo, exists := config.GetRepository("repo1")
		if !exists {
			t.Fatal("Repository should exist")
		}
		if repo.Name != "repo1" {
			t.Errorf("Name = %s, want %s", repo.Name, "repo1")
		}
		if repo.Path != "/path/to/repo1" {
			t.Errorf("Path = %s, want %s", repo.Path, "/path/to/repo1")
		}
	})

	t.Run("non-existing repository", func(t *testing.T) {
		_, exists := config.GetRepository("nonexistent")
		if exists {
			t.Error("Repository should not exist")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		emptyConfig := &Config{}
		_, exists := emptyConfig.GetRepository("repo1")
		if exists {
			t.Error("Repository should not exist in empty config")
		}
	})
}

func TestConfig_GetRepositoriesForGroup(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1", "repo2"})
	group2 := entities.NewGroup("group2", []string{"repo3", "nonexistent"})

	config := &Config{
		Repositories: map[string]*RepositoryConfig{
			"repo1": {Path: "/path/to/repo1"},
			"repo2": {Path: "/path/to/repo2"},
			"repo3": {Path: "/path/to/repo3"},
		},
		Groups: map[string]*entities.Group{
			"group1": group1,
			"group2": group2,
		},
	}

	t.Run("existing group with all repositories", func(t *testing.T) {
		repos, err := config.GetRepositoriesForGroup("group1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(repos) != 2 {
			t.Errorf("Expected 2 repositories, got %d", len(repos))
		}
		// Check if both repos are present
		repoNames := make(map[string]bool)
		for _, repo := range repos {
			repoNames[repo.Name] = true
		}
		if !repoNames["repo1"] || !repoNames["repo2"] {
			t.Error("Expected repo1 and repo2 to be present")
		}
	})

	t.Run("existing group with missing repository", func(t *testing.T) {
		repos, err := config.GetRepositoriesForGroup("group2")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Should only return repo3, skip nonexistent
		if len(repos) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(repos))
		}
		if repos[0].Name != "repo3" {
			t.Errorf("Expected repo3, got %s", repos[0].Name)
		}
	})

	t.Run("non-existing group", func(t *testing.T) {
		_, err := config.GetRepositoriesForGroup("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existing group")
		}
		var groupNotFoundErr ErrGroupNotFound
		if _, ok := err.(ErrGroupNotFound); !ok {
			t.Errorf("Expected ErrGroupNotFound, got %T", err)
		}
		groupNotFoundErr = err.(ErrGroupNotFound)
		if groupNotFoundErr.GroupName != "nonexistent" {
			t.Errorf("Expected group name 'nonexistent', got %s", groupNotFoundErr.GroupName)
		}
	})
}

func TestConfig_GetAllRepositories(t *testing.T) {
	t.Run("config with repositories", func(t *testing.T) {
		config := &Config{
			Repositories: map[string]*RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
				"repo2": {Path: "/path/to/repo2"},
				"repo3": {Path: "/path/to/repo3"},
			},
		}

		repos := config.GetAllRepositories()
		if len(repos) != 3 {
			t.Errorf("Expected 3 repositories, got %d", len(repos))
		}

		// Check all repositories are present
		repoNames := make(map[string]bool)
		for _, repo := range repos {
			repoNames[repo.Name] = true
		}
		if !repoNames["repo1"] || !repoNames["repo2"] || !repoNames["repo3"] {
			t.Error("Not all repositories were returned")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		config := &Config{}
		repos := config.GetAllRepositories()
		if len(repos) != 0 {
			t.Errorf("Expected 0 repositories, got %d", len(repos))
		}
	})

	t.Run("nil repositories map", func(t *testing.T) {
		config := &Config{Repositories: nil}
		repos := config.GetAllRepositories()
		if len(repos) != 0 {
			t.Errorf("Expected 0 repositories, got %d", len(repos))
		}
	})
}

func TestConfig_GetAllGroups(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1"})
	group2 := entities.NewGroup("group2", []string{"repo2"})

	t.Run("config with groups", func(t *testing.T) {
		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		groups := config.GetAllGroups()
		if len(groups) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(groups))
		}

		// Check all groups are present
		groupNames := make(map[string]bool)
		for _, group := range groups {
			groupNames[group.Name] = true
		}
		if !groupNames["group1"] || !groupNames["group2"] {
			t.Error("Not all groups were returned")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		config := &Config{}
		groups := config.GetAllGroups()
		if len(groups) != 0 {
			t.Errorf("Expected 0 groups, got %d", len(groups))
		}
	})

	t.Run("nil groups map", func(t *testing.T) {
		config := &Config{Groups: nil}
		groups := config.GetAllGroups()
		if len(groups) != 0 {
			t.Errorf("Expected 0 groups, got %d", len(groups))
		}
	})
}

func TestConfig_GetGroupNames(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1"})
	group2 := entities.NewGroup("group2", []string{"repo2"})

	t.Run("config with groups", func(t *testing.T) {
		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		names := config.GetGroupNames()
		if len(names) != 2 {
			t.Errorf("Expected 2 group names, got %d", len(names))
		}

		// Check all group names are present
		nameMap := make(map[string]bool)
		for _, name := range names {
			nameMap[name] = true
		}
		if !nameMap["group1"] || !nameMap["group2"] {
			t.Error("Not all group names were returned")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		config := &Config{}
		names := config.GetGroupNames()
		if len(names) != 0 {
			t.Errorf("Expected 0 group names, got %d", len(names))
		}
	})

	t.Run("nil groups map", func(t *testing.T) {
		config := &Config{Groups: nil}
		names := config.GetGroupNames()
		if len(names) != 0 {
			t.Errorf("Expected 0 group names, got %d", len(names))
		}
	})
}

func TestConfig_AddRepository(t *testing.T) {
	t.Run("add to empty config", func(t *testing.T) {
		config := &Config{}
		config.AddRepository("repo1", "/path/to/repo1")

		if config.Repositories == nil {
			t.Fatal("Repositories map should be initialized")
		}
		if len(config.Repositories) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(config.Repositories))
		}
		if config.Repositories["repo1"].Path != "/path/to/repo1" {
			t.Errorf("Path = %s, want %s", config.Repositories["repo1"].Path, "/path/to/repo1")
		}
	})

	t.Run("add to existing config", func(t *testing.T) {
		config := &Config{
			Repositories: map[string]*RepositoryConfig{
				"existing": {Path: "/existing/path"},
			},
		}
		config.AddRepository("repo1", "/path/to/repo1")

		if len(config.Repositories) != 2 {
			t.Errorf("Expected 2 repositories, got %d", len(config.Repositories))
		}
		if config.Repositories["repo1"].Path != "/path/to/repo1" {
			t.Errorf("Path = %s, want %s", config.Repositories["repo1"].Path, "/path/to/repo1")
		}
		// Ensure existing repository is still there
		if config.Repositories["existing"].Path != "/existing/path" {
			t.Error("Existing repository should still be present")
		}
	})

	t.Run("overwrite existing repository", func(t *testing.T) {
		config := &Config{
			Repositories: map[string]*RepositoryConfig{
				"repo1": {Path: "/old/path"},
			},
		}
		config.AddRepository("repo1", "/new/path")

		if len(config.Repositories) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(config.Repositories))
		}
		if config.Repositories["repo1"].Path != "/new/path" {
			t.Errorf("Path = %s, want %s", config.Repositories["repo1"].Path, "/new/path")
		}
	})
}

func TestConfig_RemoveRepository(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1", "repo2"})
	group2 := entities.NewGroup("group2", []string{"repo1", "repo3"})

	t.Run("remove existing repository", func(t *testing.T) {
		config := &Config{
			Repositories: map[string]*RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
				"repo2": {Path: "/path/to/repo2"},
				"repo3": {Path: "/path/to/repo3"},
			},
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		config.RemoveRepository("repo1")

		// Check repository is removed
		if _, exists := config.Repositories["repo1"]; exists {
			t.Error("Repository should be removed from repositories map")
		}
		if len(config.Repositories) != 2 {
			t.Errorf("Expected 2 repositories remaining, got %d", len(config.Repositories))
		}

		// Check repository is removed from groups
		if config.Groups["group1"].ContainsRepository("repo1") {
			t.Error("Repository should be removed from group1")
		}
		if config.Groups["group2"].ContainsRepository("repo1") {
			t.Error("Repository should be removed from group2")
		}

		// Check other repositories remain in groups
		if !config.Groups["group1"].ContainsRepository("repo2") {
			t.Error("repo2 should still be in group1")
		}
		if !config.Groups["group2"].ContainsRepository("repo3") {
			t.Error("repo3 should still be in group2")
		}
	})

	t.Run("remove non-existing repository", func(t *testing.T) {
		config := &Config{
			Repositories: map[string]*RepositoryConfig{
				"repo1": {Path: "/path/to/repo1"},
			},
		}

		// This should not panic
		config.RemoveRepository("nonexistent")

		// Check existing repository is still there
		if len(config.Repositories) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(config.Repositories))
		}
	})

	t.Run("remove from empty config", func(t *testing.T) {
		config := &Config{}
		// This should not panic
		config.RemoveRepository("repo1")
	})
}

func TestConfig_AddGroup(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1"})
	group2 := entities.NewGroup("group2", []string{"repo2"})

	t.Run("add to empty config", func(t *testing.T) {
		config := &Config{}
		config.AddGroup(group1)

		if config.Groups == nil {
			t.Fatal("Groups map should be initialized")
		}
		if len(config.Groups) != 1 {
			t.Errorf("Expected 1 group, got %d", len(config.Groups))
		}
		if config.Groups["group1"] != group1 {
			t.Error("Group should be added")
		}
	})

	t.Run("add to existing config", func(t *testing.T) {
		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": group1,
			},
		}
		config.AddGroup(group2)

		if len(config.Groups) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(config.Groups))
		}
		if config.Groups["group2"] != group2 {
			t.Error("New group should be added")
		}
		// Ensure existing group is still there
		if config.Groups["group1"] != group1 {
			t.Error("Existing group should still be present")
		}
	})

	t.Run("overwrite existing group", func(t *testing.T) {
		oldGroup := entities.NewGroup("group1", []string{"old_repo"})
		newGroup := entities.NewGroup("group1", []string{"new_repo"})

		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": oldGroup,
			},
		}
		config.AddGroup(newGroup)

		if len(config.Groups) != 1 {
			t.Errorf("Expected 1 group, got %d", len(config.Groups))
		}
		if config.Groups["group1"] != newGroup {
			t.Error("Group should be overwritten")
		}
	})
}

func TestConfig_RemoveGroup(t *testing.T) {
	group1 := entities.NewGroup("group1", []string{"repo1"})
	group2 := entities.NewGroup("group2", []string{"repo2"})

	t.Run("remove existing group", func(t *testing.T) {
		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": group1,
				"group2": group2,
			},
		}

		config.RemoveGroup("group1")

		if _, exists := config.Groups["group1"]; exists {
			t.Error("Group should be removed")
		}
		if len(config.Groups) != 1 {
			t.Errorf("Expected 1 group remaining, got %d", len(config.Groups))
		}
		// Ensure other group is still there
		if config.Groups["group2"] != group2 {
			t.Error("Other group should still be present")
		}
	})

	t.Run("remove non-existing group", func(t *testing.T) {
		config := &Config{
			Groups: map[string]*entities.Group{
				"group1": group1,
			},
		}

		// This should not panic
		config.RemoveGroup("nonexistent")

		// Check existing group is still there
		if len(config.Groups) != 1 {
			t.Errorf("Expected 1 group, got %d", len(config.Groups))
		}
	})

	t.Run("remove from empty config", func(t *testing.T) {
		config := &Config{}
		// This should not panic
		config.RemoveGroup("group1")
	})
}

func TestErrGroupNotFound(t *testing.T) {
	err := ErrGroupNotFound{GroupName: "test-group"}
	expected := "group 'test-group' not found"
	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}

func TestErrRepositoryNotFound(t *testing.T) {
	err := ErrRepositoryNotFound{RepositoryName: "test-repo"}
	expected := "repository 'test-repo' not found"
	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}

func TestRepositoryConfig_Fields(t *testing.T) {
	config := &RepositoryConfig{
		Path: "/path/to/repo",
	}

	if config.Path != "/path/to/repo" {
		t.Errorf("Path = %s, want %s", config.Path, "/path/to/repo")
	}
}

func TestConfig_Fields(t *testing.T) {
	repoConfig := &RepositoryConfig{Path: "/path/to/repo"}
	group := entities.NewGroup("group1", []string{"repo1"})

	config := &Config{
		Repositories: map[string]*RepositoryConfig{
			"repo1": repoConfig,
		},
		Groups: map[string]*entities.Group{
			"group1": group,
		},
		Theme:   "dark",
		Version: "1.0.0",
	}

	if len(config.Repositories) != 1 {
		t.Errorf("Repositories length = %d, want %d", len(config.Repositories), 1)
	}
	if config.Repositories["repo1"] != repoConfig {
		t.Error("Repository config should match")
	}
	if len(config.Groups) != 1 {
		t.Errorf("Groups length = %d, want %d", len(config.Groups), 1)
	}
	if config.Groups["group1"] != group {
		t.Error("Group should match")
	}
	if config.Theme != "dark" {
		t.Errorf("Theme = %s, want %s", config.Theme, "dark")
	}
	if config.Version != "1.0.0" {
		t.Errorf("Version = %s, want %s", config.Version, "1.0.0")
	}
}
