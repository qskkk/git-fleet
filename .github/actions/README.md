# Auto Tag Release Action

This GitHub action automates the creation of tags and releases based on PR content and commit messages.

## Features

- ✅ **Automatic version type detection** based on keywords
- 🏷️ **Automatic tag creation** with semantic version management
- 📦 **GitHub release creation** with automatic release notes
- 🔄 **PR support** with automatic PR number extraction
- 🛠️ **Multi-platform builds** with automatic binary uploads

## Version Bump Keywords

The action analyzes PR content and commit messages to determine the version bump type:

- **MAJOR** (1.0.0 → 2.0.0): `major`, `breaking`, `breaking-change`
- **MINOR** (1.0.0 → 1.1.0): `minor`, `feature`, `feat`
- **PATCH** (1.0.0 → 1.0.1): `patch`, `fix`, `bugfix`, `hotfix`

## Usage

### Basic Workflow

```yaml
name: Auto Release

on:
  push:
    branches: [main]

permissions:
  contents: write
  pull-requests: read

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Auto Tag Release
        uses: ./.github/actions/tag.yml
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          tag-prefix: "v"
          default-bump: "patch"
```

### Inputs

| Input          | Description                 | Required | Default               |
| -------------- | --------------------------- | -------- | --------------------- |
| `github-token` | GitHub token for API access | Yes      | `${{ github.token }}` |
| `tag-prefix`   | Prefix for tags             | No       | `'v'`                 |
| `default-bump` | Default version bump type   | No       | `'patch'`             |

### Outputs

| Output             | Description                                |
| ------------------ | ------------------------------------------ |
| `new-version`      | The new version that was created           |
| `previous-version` | The previous version tag                   |
| `bump-type`        | The type of version bump performed         |
| `tag-created`      | Whether a new tag was created (true/false) |
| `pr-number`        | The PR number if found in commit message   |
| `upload-url`       | The upload URL for the GitHub release      |
| `release-id`       | The GitHub release ID                      |

## PR Examples

### PR for a feature (MINOR bump)

```
feat: Add new dashboard feature

This PR adds a new dashboard with real-time metrics.
```

### PR for a fix (PATCH bump)

```
fix: Resolve login issue

Fixed authentication bug that prevented users from logging in.
```

### PR for a breaking change (MAJOR bump)

```
breaking: Refactor API endpoints

This is a breaking change that modifies the API structure.
```

## File Structure

```
.github/
├── actions/
│   └── tag.yml          # Composite action
└── workflows/
    └── release.yml      # Main workflow
```

## How It Works

1. **Analysis**: The action analyzes PRs and commits to determine the version type
2. **Calculation**: Calculates the new version based on the previous one
3. **Tagging**: Creates and pushes the new tag
4. **Release**: Creates a GitHub release with automatic notes
5. **Build**: Compiles and uploads binaries for different platforms

## Notes

- The action only works on main branches (main/master)
- It requires `contents: write` and `pull-requests: read` permissions
- Tags follow semantic versioning format (vX.Y.Z)
- Builds are generated for Linux, macOS, and Windows (amd64/arm64)
