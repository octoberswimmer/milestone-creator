package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v68/github"
)

type ReleaseEvent struct {
	Release struct {
		TagName string `json:"tag_name"`
	} `json:"release"`
}

type SemVer struct {
	Major  int
	Minor  int
	Patch  int
	Prefix string
}

func (v SemVer) String() string {
	return fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
}

func (v SemVer) Increment(part string) SemVer {
	switch part {
	case "major":
		return SemVer{Prefix: v.Prefix, Major: v.Major + 1, Minor: 0, Patch: 0}
	case "patch":
		return SemVer{Prefix: v.Prefix, Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	default:
		return SemVer{Prefix: v.Prefix, Major: v.Major, Minor: v.Minor + 1, Patch: 0}
	}
}

var semverRegex = regexp.MustCompile(`^(v?)(\d+)\.(\d+)\.(\d+)$`)

func ParseSemVer(tag string) (SemVer, error) {
	matches := semverRegex.FindStringSubmatch(tag)
	if matches == nil {
		return SemVer{}, fmt.Errorf("tag %q is not a valid semantic version", tag)
	}

	major, _ := strconv.Atoi(matches[2])
	minor, _ := strconv.Atoi(matches[3])
	patch, _ := strconv.Atoi(matches[4])

	return SemVer{
		Prefix: matches[1],
		Major:  major,
		Minor:  minor,
		Patch:  patch,
	}, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return fmt.Errorf("GITHUB_REPOSITORY is required")
	}

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		return fmt.Errorf("GITHUB_EVENT_PATH is required")
	}

	upcomingStr := os.Getenv("INPUT_UPCOMING_MILESTONES")
	if upcomingStr == "" {
		upcomingStr = "1"
	}
	upcoming, err := strconv.Atoi(upcomingStr)
	if err != nil {
		return fmt.Errorf("invalid upcoming-milestones value: %w", err)
	}

	increment := os.Getenv("INPUT_VERSION_INCREMENT")
	if increment == "" {
		increment = "minor"
	}
	if increment != "major" && increment != "minor" && increment != "patch" {
		return fmt.Errorf("version-increment must be major, minor, or patch")
	}

	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		return fmt.Errorf("failed to read event file: %w", err)
	}

	var event ReleaseEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	tagName := event.Release.TagName
	if tagName == "" {
		return fmt.Errorf("release has no tag")
	}

	version, err := ParseSemVer(tagName)
	if err != nil {
		fmt.Printf("Skipping: %v\n", err)
		return nil
	}

	owner, repoName, err := parseRepo(repo)
	if err != nil {
		return err
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	existingMilestones, err := listMilestones(ctx, client, owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to list milestones: %w", err)
	}

	currentVersionTitle := version.String()
	if info, exists := existingMilestones[currentVersionTitle]; exists {
		if info.State != "closed" {
			_, _, err := client.Issues.EditMilestone(ctx, owner, repoName, info.Number, &github.Milestone{
				State: github.Ptr("closed"),
			})
			if err != nil {
				return fmt.Errorf("failed to close milestone %s: %w", currentVersionTitle, err)
			}
			fmt.Printf("Closed milestone %s\n", currentVersionTitle)
		} else {
			fmt.Printf("Milestone %s already closed\n", currentVersionTitle)
		}
	} else {
		_, _, err := client.Issues.CreateMilestone(ctx, owner, repoName, &github.Milestone{
			Title: github.Ptr(currentVersionTitle),
			State: github.Ptr("closed"),
		})
		if err != nil {
			return fmt.Errorf("failed to create milestone %s: %w", currentVersionTitle, err)
		}
		fmt.Printf("Created and closed milestone %s\n", currentVersionTitle)
	}

	nextVersion := version
	for i := 0; i < upcoming; i++ {
		nextVersion = nextVersion.Increment(increment)
		title := nextVersion.String()
		if _, exists := existingMilestones[title]; exists {
			fmt.Printf("Milestone %s already exists\n", title)
			continue
		}

		_, _, err := client.Issues.CreateMilestone(ctx, owner, repoName, &github.Milestone{
			Title: github.Ptr(title),
		})
		if err != nil {
			return fmt.Errorf("failed to create milestone %s: %w", title, err)
		}
		fmt.Printf("Created milestone %s\n", title)
	}

	return nil
}

func parseRepo(repo string) (owner, name string, err error) {
	for i, c := range repo {
		if c == '/' {
			return repo[:i], repo[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid repository format: %s", repo)
}

type MilestoneInfo struct {
	Number int
	State  string
}

func listMilestones(ctx context.Context, client *github.Client, owner, repo string) (map[string]MilestoneInfo, error) {
	milestones := make(map[string]MilestoneInfo)

	opts := &github.MilestoneListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		ms, resp, err := client.Issues.ListMilestones(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, m := range ms {
			milestones[m.GetTitle()] = MilestoneInfo{
				Number: m.GetNumber(),
				State:  m.GetState(),
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return milestones, nil
}
