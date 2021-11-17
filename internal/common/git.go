package common

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/mjpitz/myago"
)

const gitRepositoryContextKey = myago.ContextKey("git.repository")

// GitRepository extracts the git repository from the context.
func GitRepository(ctx context.Context) *git.Repository {
	v := ctx.Value(gitRepositoryContextKey)
	if v == nil {
		return nil
	}

	return v.(*git.Repository)
}

// SetupGitRepository opens the current repository.
func SetupGitRepository(ctx context.Context) (context.Context, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// plane open appropriately finds the repo root
	repository, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, gitRepositoryContextKey, repository), nil
}
