package main

import (
	"testing"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    SemVer
		wantErr bool
	}{
		{
			name: "parses version with v prefix",
			tag:  "v1.2.3",
			want: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name: "parses version without prefix",
			tag:  "1.2.3",
			want: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name: "parses version with zeros",
			tag:  "v0.0.0",
			want: SemVer{Prefix: "v", Major: 0, Minor: 0, Patch: 0},
		},
		{
			name: "parses large version numbers",
			tag:  "v10.20.30",
			want: SemVer{Prefix: "v", Major: 10, Minor: 20, Patch: 30},
		},
		{
			name:    "rejects invalid format",
			tag:     "not-a-version",
			wantErr: true,
		},
		{
			name:    "rejects version with prerelease",
			tag:     "v1.2.3-beta",
			wantErr: true,
		},
		{
			name:    "rejects version with build metadata",
			tag:     "v1.2.3+build",
			wantErr: true,
		},
		{
			name:    "rejects partial version",
			tag:     "v1.2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSemVer(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseSemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_String(t *testing.T) {
	tests := []struct {
		name   string
		semver SemVer
		want   string
	}{
		{
			name:   "formats version with prefix",
			semver: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
			want:   "v1.2.3",
		},
		{
			name:   "formats version without prefix",
			semver: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3},
			want:   "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.semver.String(); got != tt.want {
				t.Errorf("SemVer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_Increment(t *testing.T) {
	tests := []struct {
		name   string
		semver SemVer
		part   string
		want   SemVer
	}{
		{
			name:   "increments major version and resets minor and patch",
			semver: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
			part:   "major",
			want:   SemVer{Prefix: "v", Major: 2, Minor: 0, Patch: 0},
		},
		{
			name:   "increments minor version and resets patch",
			semver: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
			part:   "minor",
			want:   SemVer{Prefix: "v", Major: 1, Minor: 3, Patch: 0},
		},
		{
			name:   "increments patch version",
			semver: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
			part:   "patch",
			want:   SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 4},
		},
		{
			name:   "defaults to minor for unknown part",
			semver: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3},
			part:   "unknown",
			want:   SemVer{Prefix: "v", Major: 1, Minor: 3, Patch: 0},
		},
		{
			name:   "preserves prefix when incrementing",
			semver: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3},
			part:   "minor",
			want:   SemVer{Prefix: "", Major: 1, Minor: 3, Patch: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.semver.Increment(tt.part); got != tt.want {
				t.Errorf("SemVer.Increment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		wantOwner string
		wantName  string
		wantErr   bool
	}{
		{
			name:      "parses valid repository",
			repo:      "owner/repo",
			wantOwner: "owner",
			wantName:  "repo",
		},
		{
			name:      "parses repository with hyphens",
			repo:      "my-org/my-repo",
			wantOwner: "my-org",
			wantName:  "my-repo",
		},
		{
			name:    "rejects repository without slash",
			repo:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, name, err := parseRepo(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("parseRepo() owner = %v, want %v", owner, tt.wantOwner)
				}
				if name != tt.wantName {
					t.Errorf("parseRepo() name = %v, want %v", name, tt.wantName)
				}
			}
		})
	}
}
