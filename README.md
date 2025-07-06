# 🚀 GitFleet

![Coverage](https://img.shields.io/badge/Coverage-81.2%25-brightgreen)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Test Build](https://github.com/qskkk/git-fleet/actions/workflows/test.yml/badge.svg)](https://github.com/qskkk/git-fleet/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/qskkk/git-fleet)](https://github.com/qskkk/git-fleet/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/qskkk/git-fleet)](https://github.com/qskkk/git-fleet)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-donate-yellow.svg)](https://coff.ee/qskkk)

**GitFleet** is a powerful command-line tool written in Go that helps developers manage multiple Git repositories from a single place. Designed for teams, DevOps engineers, and power users working across many projects, GitFleet simplifies routine Git operations across entire fleets of repositories.

Whether you're managing microservices, maintaining multiple projects, or coordinating across different teams, GitFleet provides both an intuitive interactive interface and powerful command-line operations to streamline your workflow.

---

## 📖 Table of Contents

- [🚀 Quick Demo](#-quick-demo)
- [✨ Features](#-features)
- [🛠️ Installation](#️-installation)
  - [📦 Nix Installation](docs/nix.md)
- [🔄 Updating GitFleet](#-updating-gitfleet)
- [🚀 Quick Start](#-quick-start)
- [📖 Usage](#-usage)
- [⚙️ Configuration](#️-configuration)
- [🔍 Automatic Repository Discovery](#-automatic-repository-discovery)
- [📂 Smart Navigation with Goto](#-smart-navigation-with-goto)
- [💡 Examples](#-examples)
- [🎨 Features in Detail](#-features-in-detail)
- [🤔 Why GitFleet?](#-why-gitfleet)
- [🔧 Advanced Usage](#-advanced-usage)
- [🏗️ Architecture](#️-architecture)
- [🛠️ Development](#️-development)
- [🤝 Contributing](#-contributing)
- [📝 License](#-license)
- [🙏 Acknowledgments](#-acknowledgments)
- [📞 Support](#-support)

---

## 🚀 Quick Demo

![GitFleet Demo](docs/media/demo.gif)

## ✨ Features

- 🎯 **Interactive Mode**: Beautiful terminal UI for easy repository and command selection
- **Auto-Discovery**: Automatically discover Git repositories in your workspace
- 🔄 **Bulk Operations**: Execute commands across multiple repositories in parallel using concurrent Go goroutines for maximum performance
- 🧩 **Smart Grouping**: Organize repositories by team, project, or any custom criteria
- 📂 **Smart Navigation**: Instantly navigate to any repository with intelligent fuzzy matching via the `goto` command
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

or

```bash
brew install qskkk/tap/git-fleet
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

### Option 4: Install with Nix Flakes

For NixOS users or those using Nix with flakes enabled:

```bash
# Install directly from the repository
nix profile install github:qskkk/git-fleet

# Or run without installing
nix run github:qskkk/git-fleet

# For development environment
nix develop github:qskkk/git-fleet
```

**Add to your NixOS configuration:**

```nix
# In your flake.nix inputs:
inputs.git-fleet.url = "github:qskkk/git-fleet";

# In your packages:
environment.systemPackages = [
  inputs.git-fleet.packages.${system}.default
];
```

**For Home Manager:**

```nix
# In your home.nix
home.packages = [
  inputs.git-fleet.packages.${system}.default
];
```

### Option 5: Install with Go

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

### Option 1: Automatic Discovery (Recommended)

1. **Install GitFleet** using one of the methods above
2. **Navigate to your workspace** containing Git repositories:

```bash
cd ~/workspace  # or wherever your Git repos are located
```

3. **Auto-discover repositories**:

```bash
gf config discover
```

4. **Start using GitFleet**:

```bash
gf                    # Interactive mode
gf status             # Status of all repositories
gf @backend pull      # Pull backend group repositories
```

### Option 2: Manual Configuration

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
gf @frontend pull    # Pull latest changes for frontend repos
gf @backend status   # Check status of backend repos
gf @all "commit -m 'Update docs'"  # Commit across all repos
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
# New multi-group syntax with @ prefix
gf @<group1> [@group2 ...] <command>

# Legacy single group syntax (still supported)
gf <group> <command>

# Examples
gf @frontend pull                    # Pull frontend repositories
gf @frontend @backend pull           # Pull both frontend and backend repositories
gf @api @database status             # Check status of api and database groups
gf @all "add . && commit -m 'fix'"   # Complex commands with quotes on all group
gf frontend pull                     # Legacy syntax still works
```

### Multi-Group Operations

GitFleet supports executing commands on multiple groups simultaneously using the `@` prefix:

```bash
# Execute on multiple groups
gf @frontend @backend @mobile pull    # Pull all three groups
gf @api @database status              # Check status of api and database
gf @team1 @team2 fetch                # Fetch updates for multiple teams
```

### Global Commands

These commands work across all repositories or provide system information:

```bash
gf config          # Show current configuration
gf config discover # Automatically discover Git repositories in current directory
gf config validate # Validate configuration file
gf config init     # Create default configuration
gf goto <repo>     # Get path to repository with fuzzy matching (for shell integration)
gf help            # Display help information
gf status          # Show status of all repositories
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

## 🔍 Automatic Repository Discovery

GitFleet can automatically discover Git repositories in your current directory and its subdirectories, making initial setup incredibly easy.

### Quick Setup with Discovery

```bash
# Navigate to your workspace directory
cd ~/workspace

# Let GitFleet discover all Git repositories
gf config discover

# View the discovered configuration
gf config

# Start using GitFleet immediately
gf status
```

The `config discover` command will:

- 🔍 **Scan recursively** for all Git repositories
- 📁 **Group by parent directory** for logical organization
- ⚙️ **Update configuration** automatically
- 🏷️ **Create smart groups** based on directory structure

### Discovery Example

If you have a workspace like this:

```
~/workspace/
├── frontend/
│   ├── web-app/        (git repo)
│   └── mobile-app/     (git repo)
├── backend/
│   ├── api-server/     (git repo)
│   └── auth-service/   (git repo)
└── tools/
    └── scripts/        (git repo)
```

Running `gf config discover` will automatically create:

- **Repositories**: web-app, mobile-app, api-server, auth-service, scripts
- **Groups**: frontend, backend, tools (based on parent directories)

---

## 📂 Smart Navigation with Goto

The `goto` command provides powerful shell integration for quick repository navigation.

### Basic Usage

```bash
# Get the path to a repository
gf goto web-app
# Output: /home/user/workspace/frontend/web-app

# Use with cd to navigate directly
cd $(gf goto web-app)
```

### Shell Integration Setup

Add this function to your shell config for seamless navigation:

**Bash/Zsh** (`~/.bashrc`, `~/.zshrc`):

```bash
# GitFleet goto function
goto() {
    cd $(gf goto "$1")
}
```

**Fish** (`~/.config/fish/config.fish`):

```fish
# GitFleet goto function
function goto
    cd (gf goto $argv[1])
end
```

### Smart Repository Matching

GitFleet now features intelligent fuzzy matching for repository names, making navigation even easier:

```bash
# Exact match (highest priority)
goto my-awesome-project  # Matches "my-awesome-project" exactly

# Partial/substring matching
goto awesome            # Matches "my-awesome-project" (contains "awesome")
goto project            # Matches "test-project" or "my-awesome-project"

# Prefix matching
goto test               # Matches "test-project" (starts with "test")

# Typo tolerance
goto awsome             # Matches "my-awesome-project" (handles missing 'e')
goto projct             # Matches "my-project" (handles missing 'e')

# Case insensitive
goto AWESOME            # Matches "my-awesome-project"
goto Test               # Matches "test-project"
```

### Advanced Goto Examples

### Shell Integration Examples

```bash
# Navigate to repository (with fuzzy matching)
goto web-app

# Quick navigation and command execution
goto api && npm install        # Matches "api-server"

# Check repository status after navigation
goto mobile && git status      # Matches "mobile-app"

# Open repository in VS Code
goto web && code .             # Matches "web-app"

# Multiple operations with fuzzy matching
goto backend && git pull && npm test  # Matches "backend-api"
```

---

## 💡 Examples

### Getting Started - Workspace Setup

```bash
# Navigate to your development workspace
cd ~/workspace

# Automatically discover all Git repositories
gf config discover

# Verify the discovered configuration
gf config

# Check status of all discovered repositories
gf status
```

### Development Workflow

```bash
# Start your day - check status of all projects
gf all status

# Fetch latest refs from all remotes for your repositories
gf all fetch

# Pull latest changes for your team's repositories
gf @frontend @backend pull

# Work on features, then commit changes
gf @frontend "add . && commit -m 'feat: new component'"

# Push all changes at once
gf @all push

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

### Quick Navigation and Repository Management

```bash
# Navigate quickly to any repository
cd $(gf goto web-app)

# Or use the shell function (if configured)
goto api-server

# Navigate and execute commands in one line
goto frontend-web && npm install && npm start

# Navigate to repository and open in VS Code
goto mobile-app && code .

# Check repository status after navigation
goto backend-api && git status

# Quick setup for new workspace
cd ~/new-workspace && gf config discover && gf status
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

- **Parallel Processing**: Lightning-fast execution across multiple repositories using concurrent Go goroutines - commands run simultaneously on all selected repositories for maximum performance
- **Error Isolation**: Failures in one repository don't stop others - each repository executes independently
- **Detailed Logging**: Clear success/failure reporting with real-time output from all repositories
- **Flexible Commands**: Support for any Git command or shell command across your entire fleet

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
- **Auto-Discovery**: Set up your entire workspace in seconds with automatic repository detection
- **Smart Navigation**: Jump to any repository instantly with the `goto` command
- **Consistent Workflows**: Standardize operations across projects
- **Time Savings**: Execute commands on dozens of repositories instantly using parallel processing with Go goroutines
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

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** with tests
4. **Run tests**: `make test`
5. **Commit your changes**: `git commit -m 'feat: add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. \*\*Open a Pull Request`

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
