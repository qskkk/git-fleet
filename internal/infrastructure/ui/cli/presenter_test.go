package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/qskkk/git-fleet/v2/internal/domain/entities"
	"github.com/qskkk/git-fleet/v2/internal/infrastructure/ui/styles"
)

func TestNewPresenter(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService)

	if presenter == nil {
		t.Error("NewPresenter() should not return nil")
	}

	// Check that it returns a Presenter type
	if _, ok := presenter.(*Presenter); !ok {
		t.Error("NewPresenter() should return a *Presenter")
	}
}

func TestPresenter_PresentExecutionSummary(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	// Create a test summary
	summary := entities.NewSummary()
	result1 := entities.NewExecutionResult("repo1", "git status")
	result1.MarkAsSuccess("Output from repo1", 0)
	summary.AddResult(*result1)

	result2 := entities.NewExecutionResult("repo2", "git status")
	result2.MarkAsFailed("Error from repo2", 1, "Command failed")
	summary.AddResult(*result2)

	summary.Finalize()

	output := presenter.PresentExecutionSummary(summary)

	if output == "" {
		t.Error("PresentExecutionSummary() should not return empty string")
	}

	// Should contain repository names
	if !contains(output, "repo1") {
		t.Error("PresentExecutionSummary() should contain 'repo1'")
	}

	if !contains(output, "repo2") {
		t.Error("PresentExecutionSummary() should contain 'repo2'")
	}
}

func TestPresenter_PresentStatusReport(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	// Create test repositories
	repos := []*entities.Repository{
		{
			Name:    "repo1",
			Path:    "/path/to/repo1",
			Status:  entities.StatusClean,
			IsValid: true,
			Branch:  "main",
		},
		{
			Name:    "repo2",
			Path:    "/path/to/repo2",
			Status:  entities.StatusModified,
			IsValid: true,
			Branch:  "develop",
		},
	}

	output := presenter.PresentStatusReport(repos)

	if output == "" {
		t.Error("PresentStatusReport() should not return empty string")
	}

	// Should contain repository names
	if !contains(output, "repo1") {
		t.Error("PresentStatusReport() should contain 'repo1'")
	}

	if !contains(output, "repo2") {
		t.Error("PresentStatusReport() should contain 'repo2'")
	}
}

func TestPresenter_PresentConfigInfo(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	// Create test data
	groups := []*entities.Group{
		entities.NewGroup("group1", []string{"repo1", "repo2"}),
	}

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1"},
		{Name: "repo2", Path: "/path/to/repo2"},
	}

	output := presenter.PresentConfigInfo(groups, repos)

	if output == "" {
		t.Error("PresentConfigInfo() should not return empty string")
	}

	// Should contain config information
	if !contains(output, "repo1") {
		t.Error("PresentConfigInfo() should contain 'repo1'")
	}

	if !contains(output, "group1") {
		t.Error("PresentConfigInfo() should contain 'group1'")
	}
}

func TestPresenter_PresentStatus(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	repos := []*entities.Repository{
		{Name: "repo1", Path: "/path/to/repo1", Status: entities.StatusClean},
	}

	ctx := context.Background()
	output, err := presenter.PresentStatus(ctx, repos, "")

	if err != nil {
		t.Errorf("PresentStatus() error = %v", err)
	}

	if output == "" {
		t.Error("PresentStatus() should not return empty string")
	}
}

func TestPresenter_PresentConfig(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	config := map[string]interface{}{
		"theme": "dark",
		"repositories": map[string]string{
			"repo1": "/path/to/repo1",
		},
	}

	ctx := context.Background()
	output, err := presenter.PresentConfig(ctx, config)

	if err != nil {
		t.Errorf("PresentConfig() error = %v", err)
	}

	if output == "" {
		t.Error("PresentConfig() should not return empty string")
	}
}

func TestPresenter_PresentSummary(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	summary := entities.NewSummary()
	result := entities.NewExecutionResult("repo1", "git status")
	result.MarkAsSuccess("Clean", 0)
	summary.AddResult(*result)
	summary.Finalize()

	ctx := context.Background()
	output, err := presenter.PresentSummary(ctx, summary)

	if err != nil {
		t.Errorf("PresentSummary() error = %v", err)
	}

	if output == "" {
		t.Error("PresentSummary() should not return empty string")
	}
}

func TestPresenter_PresentError(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	testErr := errors.New("test error message")
	ctx := context.Background()
	output := presenter.PresentError(ctx, testErr)

	if output == "" {
		t.Error("PresentError() should not return empty string")
	}

	if !contains(output, "test error message") {
		t.Error("PresentError() should contain error message")
	}
}

func TestPresenter_PresentHelp(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	ctx := context.Background()
	output := presenter.PresentHelp(ctx)

	if output == "" {
		t.Error("PresentHelp() should not return empty string")
	}

	// Should contain help information
	if !contains(output, "git-fleet") {
		t.Error("PresentHelp() should contain 'git-fleet'")
	}
}

func TestPresenter_PresentVersion(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := NewPresenter(stylesService).(*Presenter)

	ctx := context.Background()
	output := presenter.PresentVersion(ctx)

	if output == "" {
		t.Error("PresentVersion() should not return empty string")
	}

	// Should contain some version information (less strict check)
	if len(output) < 5 {
		t.Error("PresentVersion() should return meaningful version information")
	}
}

func TestPresenter_Fields(t *testing.T) {
	stylesService := styles.NewService("fleet")
	presenter := &Presenter{
		styles: stylesService,
	}

	if presenter.styles == nil {
		t.Error("Presenter styles field should not be nil")
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr ||
			len(s) >= len(substr) &&
				findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
