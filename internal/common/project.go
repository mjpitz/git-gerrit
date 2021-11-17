package common

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/mjpitz/myago"
)

const gerritProjectContextKey = myago.ContextKey("gerrit.project")

// ProjectID defines an id for a project in gerrit.
type ProjectID string

func (p ProjectID) String() string {
	return string(p)
}

// Project returns the current ProjectID
func Project(ctx context.Context) ProjectID {
	v := ctx.Value(gerritProjectContextKey)
	if v == nil {
		return ""
	}

	return v.(ProjectID)
}

// SetupProject detects and configures the current ProjectID
func SetupProject(ctx context.Context) (context.Context, error) {
	repository := GitRepository(ctx)

	originRemote, err := repository.Remote("origin")
	if errors.Is(err, git.ErrRemoteNotFound) {
		return nil, errors.New("failed to determine remote \"origin\"")
	} else if err != nil {
		return nil, err
	}

	projectID := originRemote.Config().URLs[0]
	if strings.HasPrefix(projectID, "git@") {
		projectID = projectID[strings.LastIndex(projectID, ":")+1:]
		projectID = strings.TrimSuffix(projectID, ".git")
	} else {
		u, err := url.Parse(projectID)
		if err != nil {
			return nil, err
		}

		projectID = strings.TrimPrefix(u.Path, "/")
	}

	return context.WithValue(ctx, gerritProjectContextKey, ProjectID(projectID)), nil
}
