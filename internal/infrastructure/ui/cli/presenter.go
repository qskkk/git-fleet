package cli

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/qskkk/git-fleet/internal/application/ports/output"
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
)

// Presenter implements the OutputPort interface for CLI
type Presenter struct {
	styles styles.Service
}

// NewPresenter creates a new CLI presenter
func NewPresenter(styles styles.Service) output.PresenterPort {
	return &Presenter{
		styles: styles,
	}
}

// PresentExecutionSummary presents the execution summary
func (p *Presenter) PresentExecutionSummary(summary *entities.Summary) string {
	var result bytes.Buffer

	// Title
	result.WriteString(p.styles.GetTitleStyle().Render("🚀 Execution Summary") + "\n\n")

	// Results table
	if len(summary.Results) > 0 {
		headers := []string{"Repository", "Status", "Duration", "Output"}
		rows := make([][]string, 0, len(summary.Results))

		for _, res := range summary.Results {
			status := "✅ Success"
			if res.IsFailed() {
				status = "❌ Failed"
			} else if res.IsCancelled() {
				status = "⏹️ Cancelled"
			} else if res.IsTimeout() {
				status = "⏱️ Timeout"
			}

			output := res.Output
			if len(output) > 50 {
				output = output[:47] + "..."
			}
			if output == "" && res.IsFailed() {
				output = res.ErrorMessage
				if len(output) > 50 {
					output = output[:47] + "..."
				}
			}

			duration := res.Duration.String()
			if res.Duration == 0 {
				duration = "N/A"
			}

			rows = append(rows, []string{
				res.Repository,
				status,
				duration,
				output,
			})
		}

		// Use responsive table for execution results
		tableOutput := p.styles.CreateResponsiveTable(headers, rows)
		result.WriteString(tableOutput + "\n")
	}

	// Summary statistics
	result.WriteString(p.styles.GetSectionStyle().Render("📊 Statistics:") + "\n")
	summaryData := [][]string{
		{"Total Repositories", strconv.Itoa(summary.TotalCount())},
		{"Successful", strconv.Itoa(summary.SuccessfulCount())},
		{"Failed", strconv.Itoa(summary.FailedCount())},
		{"Cancelled", strconv.Itoa(summary.CancelledCount())},
		{"Duration", summary.GetTotalDuration().String()},
	}

	statisticsHeaders := []string{"Metric", "Value"}
	statisticsTable := p.styles.CreateResponsiveTable(statisticsHeaders, summaryData)
	result.WriteString(statisticsTable + "\n")

	return result.String()
}

// PresentStatusReport presents the status report
func (p *Presenter) PresentStatusReport(repos []*entities.Repository) string {
	var result bytes.Buffer

	// Title
	result.WriteString(p.styles.GetTitleStyle().Render("📊 Repository Status Report") + "\n\n")

	if len(repos) == 0 {
		result.WriteString(p.styles.GetErrorStyle().Render("No repositories found") + "\n")
		return result.String()
	}

	// Status table
	headers := []string{"Repository", "Branch", "Status", "Changes", "Path"}
	rows := make([][]string, 0, len(repos))

	totalRepos := len(repos)
	cleanRepos := 0
	modifiedRepos := 0

	for _, repo := range repos {
		status := "✅ Clean"
		changes := "None"

		if repo.Status == "error" {
			status = "❌ Error"
			changes = "N/A"
		} else if repo.HasChanges() {
			status = "📝 Modified"
			modifiedRepos++
			
			var changesParts []string
			if repo.CreatedFiles > 0 {
				changesParts = append(changesParts, fmt.Sprintf("+%d", repo.CreatedFiles))
			}
			if repo.ModifiedFiles > 0 {
				changesParts = append(changesParts, fmt.Sprintf("~%d", repo.ModifiedFiles))
			}
			if repo.DeletedFiles > 0 {
				changesParts = append(changesParts, fmt.Sprintf("-%d", repo.DeletedFiles))
			}
			changes = strings.Join(changesParts, " ")
		} else {
			cleanRepos++
		}

		branch := repo.Branch
		if branch == "" {
			branch = "unknown"
		}

		// Use full path - let styles service handle truncation for display
		rows = append(rows, []string{
			repo.Name,
			branch,
			status,
			changes,
			repo.Path, // Use full path here
		})
	}

	// Use responsive table creation
	tableOutput := p.styles.CreateResponsiveTable(headers, rows)
	result.WriteString(tableOutput + "\n")

	// Summary statistics
	result.WriteString(p.styles.GetSectionStyle().Render("📊 Summary:") + "\n")
	summaryData := [][]string{
		{"Total Repositories", strconv.Itoa(totalRepos)},
		{"Clean Repositories", strconv.Itoa(cleanRepos)},
		{"Modified Repositories", strconv.Itoa(modifiedRepos)},
	}

	summaryHeaders := []string{"Metric", "Count"}
	summaryTable := p.styles.CreateResponsiveTable(summaryHeaders, summaryData)
	result.WriteString(summaryTable + "\n")

	return result.String()
}

// PresentConfigInfo presents configuration information
func (p *Presenter) PresentConfigInfo(groups []*entities.Group, repos []*entities.Repository) string {
	var result bytes.Buffer

	// Title
	result.WriteString(p.styles.GetTitleStyle().Render("⚙️ Configuration Information") + "\n\n")

	// Repositories section
	if len(repos) > 0 {
		result.WriteString(p.styles.GetSectionStyle().Render("📚 Repositories:") + "\n")
		
		headers := []string{"Name", "Path", "Status"}
		rows := make([][]string, 0, len(repos))

		for _, repo := range repos {
			status := "✅ Valid"
			if repo.Status == "error" {
				status = "❌ Invalid"
			}

			// Use full path - let styles service handle truncation for display
			rows = append(rows, []string{
				repo.Name,
				repo.Path, // Use full path here
				status,
			})
		}

		// Use responsive table for repositories
		repoTableOutput := p.styles.CreateResponsiveTable(headers, rows)
		result.WriteString(repoTableOutput + "\n")
	}

	// Groups section
	if len(groups) > 0 {
		result.WriteString(p.styles.GetSectionStyle().Render("🏷️ Groups:") + "\n")
		
		headers := []string{"Group", "Repositories", "Status"}
		rows := make([][]string, 0, len(groups))

		for _, group := range groups {
			status := "✅ Valid"
			repoNames := strings.Join(group.Repositories, ", ")
			
			if len(repoNames) > 50 {
				repoNames = repoNames[:47] + "..."
			}

			rows = append(rows, []string{
				group.Name,
				repoNames,
				status,
			})
		}

		// Use responsive table for groups
		groupTableOutput := p.styles.CreateResponsiveTable(headers, rows)
		result.WriteString(groupTableOutput + "\n")
	}

	return result.String()
}

// PresentStatus presents repository status information
func (p *Presenter) PresentStatus(ctx context.Context, repos []*entities.Repository, groupFilter string) (string, error) {
	var result bytes.Buffer

	// Title
	title := "📊 Repository Status Report"
	if groupFilter != "" {
		title = fmt.Sprintf("📊 Repository Status Report - Group: %s", groupFilter)
	}
	result.WriteString(p.styles.GetTitleStyle().Render(title) + "\n\n")

	if len(repos) == 0 {
		result.WriteString(p.styles.GetErrorStyle().Render("No repositories found") + "\n")
		return result.String(), nil
	}

	// Use existing PresentStatusReport logic but return with error
	statusReport := p.PresentStatusReport(repos)
	return statusReport, nil
}

// PresentConfig presents configuration information
func (p *Presenter) PresentConfig(ctx context.Context, config interface{}) (string, error) {
	var result bytes.Buffer

	// Title
	result.WriteString(p.styles.GetTitleStyle().Render("⚙️ Configuration Information") + "\n\n")
	
	// For now, just display the basic config info
	result.WriteString(p.styles.GetSectionStyle().Render("📁 Config File:") + "\n")
	result.WriteString("Configuration loaded successfully\n\n")

	return result.String(), nil
}

// PresentSummary presents execution summary
func (p *Presenter) PresentSummary(ctx context.Context, summary *entities.Summary) (string, error) {
	summaryStr := p.PresentExecutionSummary(summary)
	return summaryStr, nil
}

// PresentError presents error information
func (p *Presenter) PresentError(ctx context.Context, err error) string {
	return p.styles.GetErrorStyle().Render("❌ Error: " + err.Error())
}

// PresentHelp presents help information
func (p *Presenter) PresentHelp(ctx context.Context) string {
	help := `🚀 GitFleet - Multi-Repository Git Command Tool

USAGE:
  gf                                    # Interactive mode
  gf @<group1> [@group2 ...] <command>  # Execute on multiple groups
  gf <group> <command>                  # Execute on single group (legacy)
  gf <global-command>                   # Execute global command

GLOBAL COMMANDS:
  help, -h, --help          Show this help message
  version, -v, --version    Show version information
  config, -c, --config      Show configuration
  status, -s, --status      Show status of all repositories

GROUP COMMANDS:
  status, ls                Show status of group repositories
  pull                      Pull latest changes
  fetch                     Fetch all remotes
  <git-command>             Execute any git command

EXAMPLES:
  gf                        # Interactive mode
  gf @frontend pull         # Pull frontend repositories
  gf @api @web status       # Status of api and web groups
  gf backend "commit -m 'fix'"  # Commit with message
  gf all fetch              # Fetch all repositories

For more information, visit: https://github.com/qskkk/git-fleet
`
	return p.styles.GetHighlightStyle().Render(help)
}

// PresentVersion presents version information
func (p *Presenter) PresentVersion(ctx context.Context) string {
	return p.styles.GetHighlightStyle().Render("GitFleet v1.0.0")
}

