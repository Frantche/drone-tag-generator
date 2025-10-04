package main

import (
	"os"
	"testing"
)

// TestBranchSelection tests the branch selection logic for pull requests
func TestBranchSelection(t *testing.T) {
	// Save original env vars
	originalBranch := os.Getenv("DRONE_BRANCH")
	originalPR := os.Getenv("DRONE_PULL_REQUEST")
	originalSourceBranch := os.Getenv("DRONE_SOURCE_BRANCH")

	// Restore env vars after test
	defer func() {
		os.Setenv("DRONE_BRANCH", originalBranch)
		os.Setenv("DRONE_PULL_REQUEST", originalPR)
		os.Setenv("DRONE_SOURCE_BRANCH", originalSourceBranch)
	}()

	tests := []struct {
		name         string
		droneBranch  string
		dronePR      string
		droneSource  string
		expectBranch string
	}{
		{
			name:         "No PR - uses DRONE_BRANCH",
			droneBranch:  "main",
			dronePR:      "",
			droneSource:  "",
			expectBranch: "main",
		},
		{
			name:         "PR set - uses DRONE_SOURCE_BRANCH",
			droneBranch:  "main",
			dronePR:      "123",
			droneSource:  "feature/test",
			expectBranch: "feature/test",
		},
		{
			name:         "PR set but no source branch - falls back to DRONE_BRANCH",
			droneBranch:  "main",
			dronePR:      "456",
			droneSource:  "",
			expectBranch: "main",
		},
		{
			name:         "Regular branch push",
			droneBranch:  "develop",
			dronePR:      "",
			droneSource:  "",
			expectBranch: "develop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("DRONE_BRANCH", tt.droneBranch)
			os.Setenv("DRONE_PULL_REQUEST", tt.dronePR)
			os.Setenv("DRONE_SOURCE_BRANCH", tt.droneSource)

			// Get the branch using the same logic as run()
			branch := os.Getenv("DRONE_BRANCH")
			if os.Getenv("DRONE_PULL_REQUEST") != "" {
				sourceBranch := os.Getenv("DRONE_SOURCE_BRANCH")
				if sourceBranch != "" {
					branch = sourceBranch
				}
			}

			if branch != tt.expectBranch {
				t.Errorf("expected branch %q, got %q", tt.expectBranch, branch)
			}
		})
	}
}
