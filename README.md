# Milestone Creator

A GitHub Action that automatically creates milestones when a release is created from a semantic version tag.

## Features

- Closes the milestone matching the release tag version (creates it if it doesn't exist)
- Creates upcoming milestones with incremented versions
- Supports `v` prefix (e.g., `v1.2.3`) or plain versions (e.g., `1.2.3`)
- Supports pre-release identifiers (e.g., `v1.0.0-beta.2`); upcoming milestones increment the trailing pre-release number
- Skips upcoming milestone creation if they already exist

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

When release `v1.0.0` is created, this closes milestone `v1.0.0` and creates milestones: `v1.1.0`, `v1.2.0`, `v1.3.0`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    upcoming-milestones: '3'
```

### Use patch version increment

When release `v1.0.0` is created, this closes milestone `v1.0.0` and creates milestone `v1.0.1`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    version-increment: 'patch'
```

### Use major version increment

When release `v1.0.0` is created, this closes milestone `v1.0.0` and creates milestone `v2.0.0`

```yaml
- uses: octoberswimmer/milestone-creator@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    version-increment: 'major'
```

### Pre-release tags

When a release tag includes a pre-release identifier (e.g., `v1.0.0-beta.2`), the action closes the matching milestone and creates the next pre-release milestone(s) by incrementing the trailing number of the pre-release identifier. The `version-increment` setting is ignored for pre-releases.

For example, when release `v1.0.0-beta.2` is created with `upcoming-milestones: 2`, this closes milestone `v1.0.0-beta.2` and creates milestones `v1.0.0-beta.3` and `v1.0.0-beta.4`.
