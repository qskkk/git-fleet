package entities

import (
	"fmt"
	"slices"

	"github.com/qskkk/git-fleet/internal/pkg/errors"
)

// Group represents a logical grouping of repositories
type Group struct {
	Name         string   `json:"name"`
	Repositories []string `json:"repositories"`
	Description  string   `json:"description,omitempty"`
}

// NewGroup creates a new group with the given name and repositories
func NewGroup(name string, repositories []string) *Group {
	return &Group{
		Name:         name,
		Repositories: repositories,
	}
}

// AddRepository adds a repository to the group if it doesn't already exist
func (g *Group) AddRepository(repoName string) {
	if !g.ContainsRepository(repoName) {
		g.Repositories = append(g.Repositories, repoName)
	}
}

// RemoveRepository removes a repository from the group
func (g *Group) RemoveRepository(repoName string) {
	for i, repo := range g.Repositories {
		if repo == repoName {
			g.Repositories = append(g.Repositories[:i], g.Repositories[i+1:]...)
			break
		}
	}
}

// ContainsRepository checks if a repository is part of this group
func (g *Group) ContainsRepository(repoName string) bool {
	return slices.Contains(g.Repositories, repoName)
}

// IsEmpty returns true if the group has no repositories
func (g *Group) IsEmpty() bool {
	return len(g.Repositories) == 0
}

// Count returns the number of repositories in the group
func (g *Group) Count() int {
	return len(g.Repositories)
}

// Validate checks if the group is valid
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.ErrGroupNameEmpty
	}
	if g.IsEmpty() {
		return errors.ErrGroupMustHaveRepositories
	}
	return nil
}

// String returns a string representation of the group
func (g *Group) String() string {
	return fmt.Sprintf("Group{Name: %s, Repositories: %v}", g.Name, g.Repositories)
}
