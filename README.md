# Milestone Creator

A GitHub Action that automatically creates milestones when a release is created from a semantic version tag.

## Features

- Creates a milestone matching the release tag version
- Creates upcoming milestones with incremented versions
- Supports `v` prefix (e.g., `v1.2.3`) or plain versions (e.g., `1.2.3`)
- Skips milestone creation if they already exist

## Usage

```yaml
name: Create Milestones

on:
  release:
    types: [created]

jobs:
  create-milestones:
    runs-on: ubuntu-latest
    steps:
      - uses: octoberswimmer/milestone-creator@v1
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `github-token` | GitHub token for API access | Yes | - |
| `upcoming-milestones` | Number of upcoming milestones to create | No | `1` |
| `version-increment` | Which version part to increment (`major`, `minor`, `patch`) | No | `minor` |

## Examples

### Create 3 upcoming milestones with minor version increment

When release `v1.0.0` is created, this creates milestones: `v1.0.0`, `v1.1.0`, `v1.2.0`, `v1.3.0`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    upcoming-milestones: '3'
```

### Use patch version increment

When release `v1.0.0` is created, this creates milestones: `v1.0.0`, `v1.0.1`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    version-increment: 'patch'
```

### Use major version increment

When release `v1.0.0` is created, this creates milestones: `v1.0.0`, `v2.0.0`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    version-increment: 'major'
```
