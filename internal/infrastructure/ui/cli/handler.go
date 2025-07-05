package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/qskkk/git-fleet/internal/application/usecases"
	"github.com/qskkk/git-fleet/internal/domain/entities"
	"github.com/qskkk/git-fleet/internal/infrastructure/ui/styles"
	"github.com/qskkk/git-fleet/internal/pkg/errors"
)

// Handler handles CLI operations
type Handler struct {
	executeCommandUC *usecases.ExecuteCommandUseCase
	statusReportUC   *usecases.StatusReportUseCase
	manageConfigUC   usecases.ManageConfigUCI
	stylesService    styles.Service
}

// NewHandler creates a new CLI handler
func NewHandler(
	executeCommandUC *usecases.ExecuteCommandUseCase,
	statusReportUC *usecases.StatusReportUseCase,
	manageConfigUC usecases.ManageConfigUCI,
	stylesService styles.Service,
) *Handler {
	return &Handler{
		executeCommandUC: executeCommandUC,
		statusReportUC:   statusReportUC,
		manageConfigUC:   manageConfigUC,
		stylesService:    stylesService,
	}
}

// Execute executes a CLI command
func (h *Handler) Execute(ctx context.Context, args []string) error {
	// Parse command line arguments
	command, err := h.parseCommand(args[1:])
	if err != nil {
		return errors.WrapCommandParsingError(err)
	}

	// Handle different command types
	switch command.Type {
	case "config":
		return h.handleConfig(ctx, command.Args)
	case "status":
		return h.handleStatus(ctx, command.Groups)
	case "goto":
		return h.handleGoto(ctx, command.Args)
	case "add-repository":
		return h.handleAddRepository(ctx, command.Args)
	case "add-group":
		return h.handleAddGroup(ctx, command.Args)
	case "remove-repository":
		return h.handleRemoveRepository(ctx, command.Args)
	case "remove-group":
		return h.handleRemoveGroup(ctx, command.Args)
	case "execute":
		return h.handleExecute(ctx, command)
	default:
		return errors.WrapUnknownCommandType(command.Type)
	}
}

// Command represents a parsed CLI command
type Command struct {
	Type     string
	Groups   []string
	Args     []string
	Parallel bool
}

// parseCommand parses command line arguments
func (h *Handler) parseCommand(args []string) (*Command, error) {
	if len(args) == 0 {
		return &Command{Type: "help"}, nil
	}

	// Filter out verbose/debug flags from arguments
	filteredArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "-v" && arg != "--verbose" && arg != "-d" && arg != "--debug" {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	cmd := &Command{
		Parallel: true, // Default to parallel execution
	}

	// Check for global commands first
	switch filteredArgs[0] {
	case "help", "-h", "--help":
		cmd.Type = "help"
		return cmd, nil
	case "version", "--version":
		cmd.Type = "version"
		return cmd, nil
	case "config", "-c", "--config":
		cmd.Type = "config"
		if len(filteredArgs) > 1 {
			cmd.Args = filteredArgs[1:]
		}
		return cmd, nil
	case "status", "-s", "--status":
		cmd.Type = "status"
		if len(filteredArgs) > 1 {
			cmd.Groups = h.parseGroups(filteredArgs[1:])
		}
		return cmd, nil
	case "goto":
		cmd.Type = "goto"
		if len(filteredArgs) > 1 {
			cmd.Args = filteredArgs[1:]
		}
		return cmd, nil
	case "add":
		if len(filteredArgs) < 2 {
			return nil, errors.ErrAddCommandRequiresSubcmd
		}
		switch filteredArgs[1] {
		case "repository", "repo":
			cmd.Type = "add-repository"
			cmd.Args = filteredArgs[2:]
		case "group":
			cmd.Type = "add-group"
			cmd.Args = filteredArgs[2:]
		default:
			return nil, errors.WrapUnknownAddSubcommand(filteredArgs[1])
		}
		return cmd, nil
	case "remove", "rm":
		if len(filteredArgs) < 2 {
			return nil, errors.ErrRemoveCommandRequiresSubcmd
		}
		switch filteredArgs[1] {
		case "repository", "repo":
			cmd.Type = "remove-repository"
			cmd.Args = filteredArgs[2:]
		case "group":
			cmd.Type = "remove-group"
			cmd.Args = filteredArgs[2:]
		default:
			return nil, errors.WrapUnknownRemoveSubcommand(filteredArgs[1])
		}
		return cmd, nil
	}

	// Parse group-based commands
	i := 0
	groups := []string{}

	// Parse groups (with @ prefix) or single group
	for i < len(filteredArgs) {
		arg := filteredArgs[i]
		if strings.HasPrefix(arg, "@") {
			// Multi-group syntax: @group1 @group2 command
			groups = append(groups, strings.TrimPrefix(arg, "@"))
		} else if i == 0 && !strings.HasPrefix(arg, "-") {
			// Legacy single group syntax: group command
			groups = append(groups, arg)
		} else {
			// Start of command
			break
		}
		i++
	}

	if len(groups) == 0 {
		return nil, errors.ErrNoGroupsSpecified
	}

	if i >= len(filteredArgs) {
		return nil, errors.ErrNoCommandSpecified
	}

	// Parse command arguments
	cmdArgs := filteredArgs[i:]

	// Special handling for built-in commands
	if len(cmdArgs) == 1 {
		switch cmdArgs[0] {
		case "status", "ls":
			cmd.Type = "status"
			cmd.Groups = groups
			return cmd, nil
		}
	}

	// Regular command execution
	cmd.Type = "execute"
	cmd.Groups = groups
	cmd.Args = cmdArgs

	return cmd, nil
}

// parseGroups parses group arguments
func (h *Handler) parseGroups(args []string) []string {
	var groups []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			groups = append(groups, strings.TrimPrefix(arg, "@"))
		} else {
			groups = append(groups, arg)
		}
	}
	return groups
}

// handleConfig handles configuration commands
func (h *Handler) handleConfig(ctx context.Context, args []string) error {
	// Handle config subcommands
	if len(args) > 0 {
		switch args[0] {
		case "validate":
			return h.manageConfigUC.ValidateConfig(ctx)
		case "init", "create":
			return h.manageConfigUC.CreateDefaultConfig(ctx)
		case "discover":
			return h.manageConfigUC.DiscoverRepositories(ctx)
		default:
			return errors.WrapUnknownConfigSubcommand(args[0])
		}
	}

	// Default behavior: show config
	request := &usecases.ShowConfigInput{
		ShowGroups:       true,
		ShowRepositories: true,
		ShowValidation:   false,
	}

	response, err := h.manageConfigUC.ShowConfig(ctx, request)
	if err != nil {
		return err
	}

	fmt.Print(response.FormattedOutput)
	return nil
}

// handleStatus handles status commands
func (h *Handler) handleStatus(ctx context.Context, groups []string) error {
	request := &usecases.StatusReportInput{
		Groups: groups,
	}

	response, err := h.statusReportUC.GetStatus(ctx, request)
	if err != nil {
		return err
	}

	fmt.Print(response.FormattedOutput)
	return nil
}

// handleExecute handles command execution
func (h *Handler) handleExecute(ctx context.Context, command *Command) error {
	// Create command string from args
	commandStr := strings.Join(command.Args, " ")

	request := &usecases.ExecuteCommandInput{
		Groups:       command.Groups,
		CommandStr:   commandStr,
		Parallel:     command.Parallel,
		AllowFailure: false,
	}

	_, err := h.executeCommandUC.Execute(ctx, request)
	if err != nil {
		return err
	}

	// The progress bar already handled the output display, so we don't need to print anything else
	return nil
}

// handleAddRepository handles adding a repository
func (h *Handler) handleAddRepository(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.ErrUsageAddRepository
	}

	input := &usecases.AddRepositoryInput{
		Name: args[0],
		Path: args[1],
	}

	if err := h.manageConfigUC.AddRepository(ctx, input); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToAddRepository, err)
	}

	fmt.Printf("✅ Repository '%s' added successfully\n", input.Name)
	return nil
}

// handleAddGroup handles adding a group
func (h *Handler) handleAddGroup(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.ErrUsageAddGroup
	}

	input := &usecases.AddGroupInput{
		Name:         args[0],
		Repositories: args[1:],
	}

	if err := h.manageConfigUC.AddGroup(ctx, input); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToAddGroup, err)
	}

	fmt.Printf("✅ Group '%s' added successfully with %d repositories\n", input.Name, len(input.Repositories))
	return nil
}

// handleRemoveRepository handles removing a repository
func (h *Handler) handleRemoveRepository(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageRemoveRepository
	}

	name := args[0]

	if err := h.manageConfigUC.RemoveRepository(ctx, name); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToRemoveRepository, err)
	}

	fmt.Printf("✅ Repository '%s' removed successfully\n", name)
	return nil
}

// handleRemoveGroup handles removing a group
func (h *Handler) handleRemoveGroup(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageRemoveGroup
	}

	name := args[0]

	if err := h.manageConfigUC.RemoveGroup(ctx, name); err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToRemoveGroup, err)
	}

	fmt.Printf("✅ Group '%s' removed successfully\n", name)
	return nil
}

// handleGoto handles the goto command to return repository paths
func (h *Handler) handleGoto(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.ErrUsageGoto
	}

	repoName := args[0]

	// Get repositories from config
	repos, err := h.manageConfigUC.GetRepositories(ctx)
	if err != nil {
		return errors.WrapRepositoryOperationError(errors.ErrFailedToGetRepositories, err)
	}

	if len(repos) == 0 {
		return errors.WrapRepositoryNotFound(repoName)
	}

	// First try exact match
	for _, repo := range repos {
		if repo.Name == repoName {
			fmt.Print(repo.Path)
			return nil
		}
	}

	// If no exact match, find the closest match
	var (
		bestMatch *entities.Repository
		bestScore float64 = 0
	)

	for _, repo := range repos {
		score := h.calculateSimilarity(repoName, repo.Name)
		if score > bestScore {
			bestScore = score
			bestMatch = repo
		}
	}

	if bestMatch == nil {
		return errors.WrapRepositoryNotFound(repoName)
	}

	// Just print the path - no styling or additional output
	fmt.Print(bestMatch.Path)
	return nil
}

// calculateSimilarity calculates the similarity between two strings
// Returns a score between 0 and 1, where 1 is identical
func (h *Handler) calculateSimilarity(a, b string) float64 {
	// Convert to lowercase for case-insensitive comparison
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	// Handle identical strings
	if a == b {
		return 1.0
	}

	// Handle empty strings
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	// Check if one string contains the other
	if strings.Contains(b, a) || strings.Contains(a, b) {
		return 0.9 // High score for substring matches
	}

	// Check if they start with the same prefix
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	// Calculate prefix similarity
	prefixMatch := 0
	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			prefixMatch++
		} else {
			break
		}
	}

	// Calculate Levenshtein distance-based similarity
	distance := h.levenshteinDistance(a, b)
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	distanceSimilarity := 1.0 - float64(distance)/float64(maxLen)
	prefixSimilarity := float64(prefixMatch) / float64(minLen)

	// Weight prefix similarity higher, and boost overall similarity for very similar strings
	similarity := 0.6*prefixSimilarity + 0.4*distanceSimilarity

	// Boost similarity for strings that differ by only a few characters
	if distance == 1 && maxLen > 3 {
		similarity = 0.85 // High similarity for single character differences
	}

	return similarity
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func (h *Handler) levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = 1 + min(
					matrix[i-1][j],   // deletion
					matrix[i][j-1],   // insertion
					matrix[i-1][j-1], // substitution
				)
			}
		}
	}

	return matrix[len(a)][len(b)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
