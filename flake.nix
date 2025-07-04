{
  description = "Git Fleet - Multi-Repository Git Command Tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        git-fleet = pkgs.buildGoModule rec {
          pname = "git-fleet";
          version = "2.5.0";

          src = ./.;

          vendorHash = "sha256-ItWBQ02MxpDWWuj56diO0MlhgaONLPvWvnf1VyzlOLU=";

          subPackages = [ "cmd/gf" ];

          ldflags = [
            "-s"
            "-w"
            "-X github.com/qskkk/git-fleet/internal/pkg/version.Version=${version}"
            "-X github.com/qskkk/git-fleet/internal/pkg/version.Commit=${src.rev or "unknown"}"
            "-X github.com/qskkk/git-fleet/internal/pkg/version.Date=1970-01-01T00:00:00Z"
          ];

          # Disable tests during build since they might require git configuration
          doCheck = false;

          meta = with pkgs.lib; {
            description = "Multi-Repository Git Command Tool for managing multiple Git repositories";
            longDescription = ''
              Git Fleet is a powerful command-line tool designed to streamline Git operations
              across multiple repositories. It allows you to organize repositories into groups
              and execute Git commands on multiple repositories simultaneously.

              Features:
              - Execute Git commands across multiple repositories
              - Organize repositories into logical groups
              - Parallel command execution for faster operations
              - Interactive and CLI modes
              - Status reporting across all repositories
              - Configurable repository and group management
            '';
            homepage = "https://github.com/qskkk/git-fleet";
            license = licenses.mit;
            maintainers = [ ];
            mainProgram = "gf";
            platforms = platforms.unix;
          };
        };
      in
      {
        packages = {
          default = git-fleet;
          git-fleet = git-fleet;
        };

        apps = {
          default = flake-utils.lib.mkApp {
            drv = git-fleet;
            name = "gf";
          };
          git-fleet = flake-utils.lib.mkApp {
            drv = git-fleet;
            name = "gf";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            git
            gnumake
            golangci-lint
            delve # Go debugger
          ];

          shellHook = ''
            echo "🚀 Git Fleet development environment"
            echo "Available commands:"
            echo "  make build       - Build the application"
            echo "  make test        - Run tests"
            echo "  make test-cover  - Run tests with coverage"
            echo "  make lint        - Run linter"
            echo "  go run cmd/gf/main.go - Run the application directly"
          '';
        };
      }
    );
}
