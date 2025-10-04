package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultMasterBranch = "main"
	tagsFileName        = ".tags"
)

func run() error {
	masterBranch := os.Getenv("PLUGIN_MASTER_BRANCH")
	if masterBranch == "" {
		masterBranch = defaultMasterBranch
	}

	branch := os.Getenv("DRONE_BRANCH")
	if branch == "" {
		return fmt.Errorf("DRONE_BRANCH non défini")
	}
	// Use DRONE_SOURCE_BRANCH for pull requests
	if os.Getenv("DRONE_PULL_REQUEST") != "" {
		sourceBranch := os.Getenv("DRONE_SOURCE_BRANCH")
		if sourceBranch != "" {
			branch = sourceBranch
		}
	}
	build := os.Getenv("DRONE_BUILD_NUMBER")
	if build == "" {
		build = "0"
	}

	nextVersion, hasNextVersion, err := GetNextVersionFromGit()
	if err != nil {
		return fmt.Errorf("échec get-next-version: %w", err)
	}

	newVersion, err := bumpVersion(nextVersion.String(), branch == masterBranch, hasNextVersion, branch, build)
	if err != nil {
		return fmt.Errorf("calcul de version: %w", err)
	}

	content := newVersion
	if branch == masterBranch {
		content = fmt.Sprintf("%s,latest", newVersion)
	}

	tmp := filepath.Join(".", tagsFileName+".tmp")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return fmt.Errorf("écriture fichier temporaire: %w", err)
	}
	if err := os.Rename(tmp, tagsFileName); err != nil {
		return fmt.Errorf("renommage vers %s: %w", tagsFileName, err)
	}

	fmt.Printf("Fichier %s généré :\n%s\n", tagsFileName, content)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Erreur: %v", err)
	}
}
