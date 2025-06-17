package entities

import (
	"time"
)

// RepositoryStatus represents the current state of a repository
type RepositoryStatus string

const (
	StatusClean    RepositoryStatus = "Clean"
	StatusModified RepositoryStatus = "Modified"
	StatusError    RepositoryStatus = "Error"
	StatusWarning  RepositoryStatus = "Warning"
	StatusCreated  RepositoryStatus = "Created"
	StatusDeleted  RepositoryStatus = "Deleted"
	StatusUnknown  RepositoryStatus = "Unknown"
)

// Repository represents a Git repository with its metadata
type Repository struct {
	Name         string           `json:"name"`
	Path         string           `json:"path"`
	Status       RepositoryStatus `json:"status"`
	Branch       string           `json:"branch"`
	CreatedFiles int              `json:"created_files"`
	ModifiedFiles int             `json:"modified_files"`
	DeletedFiles int              `json:"deleted_files"`
	LastChecked  time.Time        `json:"last_checked"`
	IsValid      bool             `json:"is_valid"`
	ErrorMessage string           `json:"error_message,omitempty"`
}

// HasChanges returns true if the repository has any pending changes
func (r *Repository) HasChanges() bool {
	return r.CreatedFiles > 0 || r.ModifiedFiles > 0 || r.DeletedFiles > 0
}

// IsHealthy returns true if the repository is in a good state
func (r *Repository) IsHealthy() bool {
	return r.IsValid && r.Status != StatusError
}

// GetDisplayPath returns a truncated path for display purposes
func (r *Repository) GetDisplayPath(maxLength int) string {
	if len(r.Path) <= maxLength {
		return r.Path
	}
	return "..." + r.Path[len(r.Path)-maxLength+3:]
}

// UpdateStatus updates the repository status based on the current state
func (r *Repository) UpdateStatus() {
	if !r.IsValid {
		r.Status = StatusError
		return
	}

	if r.HasChanges() {
		r.Status = StatusModified
		return
	}

	r.Status = StatusClean
}
