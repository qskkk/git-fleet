# Auto Tag Release Action

Cette action GitHub automatise la création de tags et de releases basée sur le contenu des PRs et des messages de commit.

## Fonctionnalités

- ✅ **Détection automatique du type de version** basée sur des mots-clés
- 🏷️ **Création automatique de tags** avec gestion sémantique des versions
- 📦 **Création de releases GitHub** avec notes de release automatiques
- 🔄 **Support des PRs** avec extraction automatique du numéro de PR
- 🛠️ **Build multi-plateforme** avec upload automatique des binaires

## Mots-clés pour le versioning

L'action analyse le contenu des PRs et les messages de commit pour déterminer le type de bump de version :

- **MAJOR** (1.0.0 → 2.0.0) : `major`, `breaking`, `breaking-change`
- **MINOR** (1.0.0 → 1.1.0) : `minor`, `feature`, `feat`
- **PATCH** (1.0.0 → 1.0.1) : `patch`, `fix`, `bugfix`, `hotfix`

## Usage

### Workflow basique

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

| Input          | Description             | Required | Default               |
| -------------- | ----------------------- | -------- | --------------------- |
| `github-token` | Token GitHub pour l'API | Oui      | `${{ github.token }}` |
| `tag-prefix`   | Préfixe pour les tags   | Non      | `'v'`                 |
| `default-bump` | Type de bump par défaut | Non      | `'patch'`             |

### Outputs

| Output             | Description                |
| ------------------ | -------------------------- |
| `new-version`      | Nouvelle version créée     |
| `previous-version` | Version précédente         |
| `bump-type`        | Type de bump effectué      |
| `tag-created`      | Tag créé (true/false)      |
| `pr-number`        | Numéro de PR si trouvé     |
| `upload-url`       | URL d'upload de la release |
| `release-id`       | ID de la release GitHub    |

## Exemples de PRs

### PR pour une feature (MINOR bump)

```
feat: Add new dashboard feature

This PR adds a new dashboard with real-time metrics.
```

### PR pour un fix (PATCH bump)

```
fix: Resolve login issue

Fixed authentication bug that prevented users from logging in.
```

### PR pour un breaking change (MAJOR bump)

```
breaking: Refactor API endpoints

This is a breaking change that modifies the API structure.
```

## Structure des fichiers

```
.github/
├── actions/
│   └── tag.yml          # Action composite
└── workflows/
    └── release.yml      # Workflow principal
```

## Fonctionnement

1. **Analyse** : L'action analyse les PRs et commits pour déterminer le type de version
2. **Calcul** : Calcule la nouvelle version basée sur la précédente
3. **Tag** : Crée et pousse le nouveau tag
4. **Release** : Crée une release GitHub avec notes automatiques
5. **Build** : Compile et upload les binaires pour différentes plateformes

## Notes

- L'action fonctionne uniquement sur les branches principales (main/master)
- Elle nécessite les permissions `contents: write` et `pull-requests: read`
- Les tags suivent le format sémantique (vX.Y.Z)
- Les builds sont générés pour Linux, macOS et Windows (amd64/arm64)
