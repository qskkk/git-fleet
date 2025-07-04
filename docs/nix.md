# Nix Installation Guide

This document explains how to install and use Git Fleet with Nix and NixOS.

## Prerequisites

- Nix package manager with flakes enabled
- For NixOS users: `nix.settings.experimental-features = [ "nix-command" "flakes" ];` in your configuration

## Installation Methods

### 1. Direct Installation

Install Git Fleet directly into your profile:

```bash
nix profile install github:qskkk/git-fleet
```

### 2. Temporary Usage

Run Git Fleet without installing:

```bash
# Run with arguments
nix run github:qskkk/git-fleet -- --help

# Run interactively
nix run github:qskkk/git-fleet
```

### 3. NixOS System Configuration

Add to your NixOS system configuration:

```nix
# flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    git-fleet.url = "github:qskkk/git-fleet";
  };

  outputs = { self, nixpkgs, git-fleet }: {
    nixosConfigurations.your-hostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux"; # or your system
      modules = [
        ({ pkgs, ... }: {
          environment.systemPackages = [
            git-fleet.packages.x86_64-linux.default
          ];
        })
      ];
    };
  };
}
```

### 4. Home Manager Configuration

For user-level installation with Home Manager:

```nix
# home.nix or flake.nix with home-manager
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    home-manager.url = "github:nix-community/home-manager";
    git-fleet.url = "github:qskkk/git-fleet";
  };

  outputs = { nixpkgs, home-manager, git-fleet, ... }: {
    homeConfigurations.your-username = home-manager.lib.homeManagerConfiguration {
      pkgs = nixpkgs.legacyPackages.x86_64-linux; # or your system
      modules = [
        ({ pkgs, ... }: {
          home.packages = [
            git-fleet.packages.x86_64-linux.default
          ];
        })
      ];
    };
  };
}
```

## Development Environment

Get a development shell with all necessary tools:

```bash
# Enter development environment
nix develop github:qskkk/git-fleet

# Or for a specific shell
nix develop github:qskkk/git-fleet#default
```

This provides:

- Go compiler
- Git
- Make
- golangci-lint
- Delve debugger

## Building from Source

Clone and build locally:

```bash
git clone https://github.com/qskkk/git-fleet.git
cd git-fleet
nix build
```

The built binary will be in `result/bin/gf`.

## Shell Integration

### For Bash/Zsh

Add to your shell configuration:

```bash
# If using nix profile
export PATH="$HOME/.nix-profile/bin:$PATH"

# Or use nix run for occasional use
alias gf='nix run github:qskkk/git-fleet --'
```

### For Fish

```fish
# If using nix profile
set -x PATH $HOME/.nix-profile/bin $PATH

# Or use nix run for occasional use
alias gf='nix run github:qskkk/git-fleet --'
```

## Updating

### Profile Installation

```bash
nix profile upgrade git-fleet
```

### System/Home Manager

Update your flake inputs:

```bash
nix flake update
# Then rebuild your system or home configuration
```

## Uninstalling

### Profile Installation

```bash
nix profile remove git-fleet
```

### System/Home Manager

Remove from your configuration and rebuild.

## Troubleshooting

### Flakes Not Enabled

If you get an error about flakes not being enabled:

```bash
# Enable flakes temporarily
nix --experimental-features 'nix-command flakes' run github:qskkk/git-fleet

# Or enable permanently in ~/.config/nix/nix.conf
echo "experimental-features = nix-command flakes" >> ~/.config/nix/nix.conf
```

### Permission Issues

If you encounter permission issues with the Nix store:

```bash
# Make sure you're in the nix-users group (macOS) or trusted-users (Linux)
sudo dscl . -append /Groups/nixbld GroupMembership $USER  # macOS
```

### Build Failures

If the build fails:

1. Make sure you have the latest nixpkgs
2. Check that your system is supported
3. Report issues at: https://github.com/qskkk/git-fleet/issues

## Advanced Usage

### Pinning to a Specific Version

```bash
# Pin to a specific commit
nix run github:qskkk/git-fleet/abc123def

# Pin to a specific tag
nix run github:qskkk/git-fleet/v1.0.0
```

### Custom Build Options

You can override build options:

```nix
# In your flake.nix
git-fleet.packages.x86_64-linux.default.override {
  buildGoModule = args: pkgs.buildGoModule (args // {
    # Custom build flags
    ldflags = args.ldflags ++ [ "-X main.customFlag=value" ];
  });
}
```

## Support

For Nix-specific issues:

- Check the flake.nix file in the repository
- Open an issue at https://github.com/qskkk/git-fleet/issues
- Tag issues with "nix" for faster response
