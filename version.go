package main

import (
    "fmt"

    gogit "github.com/go-git/go-git/v5"
    "github.com/Masterminds/semver"
    "github.com/thenativeweb/get-next-version/git"
    "github.com/thenativeweb/get-next-version/versioning"
)

func GetNextVersionFromGit() (*semver.Version, bool, error) {
    repository, err := gogit.PlainOpen(".")
    if err != nil {
        return nil, false, fmt.Errorf("échec ouverture du repo git: %w", err)
    }

    result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
    if err != nil {
        return nil, false, fmt.Errorf("échec analyse des commits: %w", err)
    }

    nextVersion, hasNextVersion := versioning.CalculateNextVersion(result.LatestReleaseVersion, result.ConventionalCommitTypes)
    return &nextVersion, hasNextVersion, nil
}

func bumpVersion(baseVersion string, isMaster, hasNextVersion bool, branch, build string) (string, error) {
    if isMaster {
        if !hasNextVersion {
            v, err := semver.NewVersion(baseVersion)
            if err != nil {
                return "", fmt.Errorf("version semver invalide: %w", err)
            }
            return fmt.Sprintf("%d.%d.%d", v.Major(), v.Minor(), v.Patch()+1), nil
        }
        return baseVersion, nil
    }
    pre := sanitizePrerelease(branch)
    return fmt.Sprintf("%s-%s.%s", baseVersion, pre, build), nil
}

func sanitizePrerelease(s string) string {
    var out []rune
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
            out = append(out, r)
        } else {
            out = append(out, '-')
        }
    }
    res := string(out)
    for len(res) > 0 && res[0] == '-' {
        res = res[1:]
    }
    for len(res) > 0 && res[len(res)-1] == '-' {
        res = res[:len(res)-1]
    }
    if res == "" {
        return "unknown"
    }
    return res
}
