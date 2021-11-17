package common

import (
	"context"
	"errors"
	"net/url"

	"github.com/go-git/go-git/v5"
	"github.com/mjpitz/myago"
	"golang.org/x/build/gerrit"
)

const gerritClientContextKey = myago.ContextKey("gerrit.client")

// GerritAPI extracts the gerrit client from the context.
func GerritAPI(ctx context.Context) *Gerrit {
	v := ctx.Value(gerritClientContextKey)
	if v == nil {
		return nil
	}

	return v.(*Gerrit)
}

// SetupGerritAPI sets up the gerrit API to be shared across multiple commands.
func SetupGerritAPI(ctx context.Context) (context.Context, error) {
	repository := GitRepository(ctx)

	gerritRemote, err := repository.Remote("gerrit")
	if errors.Is(err, git.ErrRemoteNotFound) {
		return nil, errors.New("failed to determine remote \"gerrit\"")
	} else if err != nil {
		return nil, err
	}

	gerritURL, err := url.Parse(gerritRemote.Config().URLs[0])
	if err != nil {
		return nil, err
	}

	g := &Gerrit{
		BaseURL: "https://"+gerritURL.Hostname(),
	}
	g.Client = gerrit.NewClient(g.BaseURL, gerrit.NoAuth)

	return context.WithValue(ctx, gerritClientContextKey, g), nil
}

// Gerrit provides programmatic access to the Gerrit REST API.
type Gerrit struct {
	BaseURL string
	Client *gerrit.Client
}
