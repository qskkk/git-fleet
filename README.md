# 🚀 GitFleet

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Test Build](https://github.com/qskkk/git-fleet/actions/workflows/test.yml/badge.svg)](https://github.com/qskkk/git-fleet/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/qskkk/git-fleet)](https://github.com/qskkk/git-fleet/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qskkk/git-fleet)](https://github.com/qskkk/git-fleet)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-donate-yellow.svg)](https://coff.ee/qskkk)

**GitFleet** is a powerful command-line tool written in Go that helps developers manage multiple Git repositories from a single place. Designed for teams, DevOps engineers, and power users working across many projects, GitFleet simplifies routine Git operations across entire fleets of repositories.

Whether you're managing microservices, maintaining multiple projects, or coordinating across different teams, GitFleet provides both an intuitive interactive interface and powerful command-line operations to streamline your workflow.

---

## ✨ Features

- 🎯 **Interactive Mode**: Beautiful terminal UI for easy repository and command selection
- 🔄 **Bulk Operations**: pull, fetch, and execute commands across multiple repositories
- 🧩 **Smart Grouping**: Organize repositories by team, project, or any custom criteria
- ⚙️ **Flexible Commands**: Run Git commands or any shell commands across your entire fleet
- ⚡ **Fast & Lightweight**: Written in Go for optimal performance
- 📁 **Simple Configuration**: Easy-to-manage JSON configuration file
- 📊 **Rich Status Reports**: Beautiful, colorized output with detailed repository status
- 🎨 **Modern UI**: Styled terminal interface with icons and colors

---

## 🛠️ Installation

### Option 1: Install with Homebrew (Recommended)

**macOS and Linux:**

```bash
brew tap qskkk/tap
brew install git-fleet
```

### Option 2: Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/qskkk/git-fleet/releases):

**Linux (x64):**

```bash
curl -L https://github.com/qskkk/git-fleet/releases/latest/download/git-fleet-linux-amd64.tar.gz | tar -xz
sudo mv git-fleet /usr/local/bin/
```

**Linux (ARM64):**

```bash
curl -L https://github.com/qskkk/git-fleet/releases/latest/download/git-fleet-linux-arm64.tar.gz | tar -xz
sudo mv git-fleet /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L https://github.com/qskkk/git-fleet/releases/latest/download/git-fleet-darwin-amd64.tar.gz | tar -xz
sudo mv git-fleet /usr/local/bin/
```

**macOS (Apple Silicon):**

```bash
curl -L https://github.com/qskkk/git-fleet/releases/latest/download/git-fleet-darwin-arm64.tar.gz | tar -xz
sudo mv git-fleet /usr/local/bin/
```

**Windows:**

```powershell
# Download the latest Windows release
curl -L -o git-fleet.zip https://github.com/qskkk/git-fleet/releases/latest/download/git-fleet-windows-amd64.zip
# Extract and add to PATH
```

### Option 3: Build from Source

**Prerequisites:** Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/qskkk/git-fleet.git
cd git-fleet

# Build and install
make install
```

### Option 4: Install with Go

```bash
go install github.com/qskkk/git-fleet@latest
```

---

## 🔄 Updating GitFleet

### Update with Homebrew

If you installed GitFleet using Homebrew, you can easily update to the latest version:

```bash
brew update
brew upgrade git-fleet
```

### Update from Source

If you built from source, navigate to your git-fleet directory and rebuild:

```bash
cd git-fleet
git pull origin main
make install
```

### Update with Go

If you installed with `go install`, simply run the install command again:

```bash
go install github.com/qskkk/git-fleet@latest
```

---

## 🚀 Quick Start

1. **Install GitFleet** using one of the methods above
2. **Create a configuration file** at `~/.config/git-fleet/.gfconfig.json`:

```json
{
  "repositories": {
    "web-app": {
      "path": "/path/to/your/web-app"
    },
    "api-server": {
      "path": "/path/to/your/api"
    },
    "mobile-app": {
      "path": "/path/to/your/mobile"
    }
  },
  "groups": {
    "frontend": ["web-app", "mobile-app"],
    "backend": ["api-server"],
    "all": ["web-app", "api-server", "mobile-app"]
  }
}
```

3. **Run GitFleet**:

```bash
# Interactive mode - select groups and commands via UI
gf

# Or use direct commands
gf frontend pull    # Pull latest changes for frontend repos
gf backend status   # Check status of backend repos
gf all "commit -m 'Update docs'"  # Commit across all repos
```

---

## 📖 Usage

### Interactive Mode

Simply run `gf` without arguments to enter interactive mode:

```bash
gf
```

This launches a beautiful terminal UI where you can:

- ✅ Select multiple repository groups
- 🎯 Choose commands to execute
- 📊 View execution results with rich formatting

### Command Line Mode

Execute commands directly on specific groups:

```bash
# Basic syntax
gf <group> <command>

# Examples
gf frontend pl              # Pull all frontend repositories
gf backend st            # Check status of backend repositories
gf backend fa             # Fetch all remotes for backend repositories
gf all "add . && commit -m 'fix'"  # Complex commands with quotes
```

### Global Commands

These commands work across all repositories or provide system information:

```bash
gf config    # Show current configuration
gf help      # Display help information
gf status    # Show status of all repositories
```

---

## ⚙️ Configuration

GitFleet uses a JSON configuration file located at `~/.config/git-fleet/.gfconfig.json`.

### Configuration Structure

```json
{
  "repositories": {
    "repo-name": {
      "path": "/absolute/path/to/repository"
    }
  },
  "groups": {
    "group-name": ["repo1", "repo2", "repo3"]
  }
}
```

### Example Configuration

```json
{
  "repositories": {
    "frontend-web": {
      "path": "/home/user/projects/webapp"
    },
    "frontend-mobile": {
      "path": "/home/user/projects/mobile-app"
    },
    "backend-api": {
      "path": "/home/user/projects/api-server"
    },
    "backend-auth": {
      "path": "/home/user/projects/auth-service"
    },
    "shared-components": {
      "path": "/home/user/projects/ui-components"
    },
    "documentation": {
      "path": "/home/user/projects/docs"
    }
  },
  "groups": {
    "frontend": ["frontend-web", "frontend-mobile", "shared-components"],
    "backend": ["backend-api", "backend-auth"],
    "mobile": ["frontend-mobile"],
    "web": ["frontend-web", "shared-components"],
    "docs": ["documentation"],
    "all": [
      "frontend-web",
      "frontend-mobile",
      "backend-api",
      "backend-auth",
      "shared-components",
      "documentation"
    ]
  }
}
```

### Configuration Tips

- **Absolute Paths**: Always use absolute paths for repository locations
- **Logical Grouping**: Create groups that match your workflow (by team, technology, environment)
- **Overlapping Groups**: Repositories can belong to multiple groups
- **Validation**: Use `gf config` to verify your configuration

---

## 💡 Examples

### Development Workflow

```bash
# Start your day - check status of all projects
gf all status

# Fetch latest refs from all remotes for your repositories
gf all fetch

# Pull latest changes for your team's repositories
gf frontend pull
gf backend pull

# Work on features, then commit changes
gf frontend "add . && commit -m 'feat: new component'"

# Push all changes at once
gf all push

# Check final status
gf all status
```

### Release Management

```bash
# Check status before release
gf production status

# Create release branches
gf production "checkout -b release/v1.2.0"

# Tag release
gf production "tag -a v1.2.0 -m 'Release v1.2.0'"

# Push tags
gf production "push --tags"
```

### Team Coordination

```bash
# Fetch latest refs from all remotes to see what's new
gf all fetch

# Update all repositories to latest
gf all pull

# Check for uncommitted changes across teams
gf all status

# Run tests across all services
gf backend "npm test"
gf frontend "npm run test"
```

---

## 🎨 Features in Detail

### Interactive Terminal UI

- **Multi-selection**: Use spacebar to select multiple groups
- **Keyboard Navigation**: Arrow keys for navigation, Enter to confirm
- **Visual Feedback**: Colorized output with status indicators
- **Error Handling**: Graceful handling of command failures

### Rich Status Reports

GitFleet provides detailed status information with:

- ✅ Clean repositories
- 📝 Repositories with changes
- 🆕 New files count
- ✏️ Modified files count
- 🗑️ Deleted files count
- ❌ Error indicators for invalid paths

### Command Execution

- **Parallel Processing**: Fast execution across multiple repositories
- **Error Isolation**: Failures in one repository don't stop others
- **Detailed Logging**: Clear success/failure reporting
- **Flexible Commands**: Support for any Git command or shell command

---

## 🤔 Why GitFleet?

### Problem It Solves

Managing multiple Git repositories manually is time-consuming and error-prone:

- Switching between directories to run the same command
- Forgetting to update certain repositories
- Inconsistent workflow across different projects
- No overview of the state of multiple repositories

### GitFleet Solution

- **Centralized Management**: Control all repositories from one place
- **Consistent Workflows**: Standardize operations across projects
- **Time Savings**: Execute commands on dozens of repositories instantly
- **Better Visibility**: Clear overview of all repository states
- **Reduced Errors**: Less manual work means fewer mistakes

### Perfect For

- **Microservices Architecture**: Manage multiple service repositories
- **Multi-Project Teams**: Coordinate across different projects
- **DevOps Engineers**: Automate repository maintenance tasks
- **Open Source Maintainers**: Manage multiple project repositories
- **Development Teams**: Standardize development workflows

---

## 🔧 Advanced Usage

### Complex Commands

Use quotes for complex commands:

```bash
# Multiple commands with &&
gf backend "git add . && git commit -m 'fix: critical bug' && git push"

# Commands with pipes
gf all "git log --oneline | head -5"

# Environment-specific commands
gf production "git checkout main && git pull && npm install"
```

### Conditional Execution

GitFleet continues execution even if some repositories fail:

```bash
# This will attempt to pull all repositories
# If some fail (e.g., merge conflicts), others continue
gf all pull
```

### Status Filtering

```bash
# Check status of specific group
gf frontend status

# Global status (all repositories)
gf status
```

---

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/qskkk/git-fleet.git
cd git-fleet

# Install dependencies
go mod download

# Run tests
make test

# Build
make install
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# View coverage report
go tool cover -html=coverage.out
```

### Project Structure

```
git-fleet/
├── main.go              # Application entry point
├── command/             # Command execution logic
├── config/              # Configuration management
├── interactive/         # Terminal UI components
├── style/               # UI styling and formatting
├── .github/             # GitHub Actions workflows
└── README.md           # This file
```

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** with tests
4. **Run tests**: `make test`
5. **Commit your changes**: `git commit -m 'feat: add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Guidelines

- Write tests for new features
- Follow Go conventions and best practices
- Update documentation for new features
- Use conventional commit messages
- Ensure all tests pass before submitting

---

## 📝 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the interactive terminal UI
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal output
- Uses [Charm](https://github.com/charmbracelet) libraries for enhanced CLI experience

---

## 📞 Support

- 📋 **Issues**: [GitHub Issues](https://github.com/qskkk/git-fleet/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/qskkk/git-fleet/discussions)
- 📖 **Documentation**: Run `gf help` for built-in help

---

<div align="center">
  <strong>⭐ Star this project if you find it helpful!</strong>
</div>
