package checkout

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mjpitz/git-gerrit/internal/common"
	"github.com/urfave/cli/v2"
	"golang.org/x/build/gerrit"
)

var (
	Command = &cli.Command{
		Name:      "checkout",
		UsageText: "checkout [options] <changeID>",
		Usage:     "Checkout an available change",
		Action: func(ctx *cli.Context) error {
			gerritAPI := common.GerritAPI(ctx.Context)
			repo := common.GitRepository(ctx.Context)

			wt, err := repo.Worktree()
			if err != nil {
				return fmt.Errorf("failed to obtain worktree: %v", err)
			}

			changeID := ctx.Args().Get(0)

			change, err := gerritAPI.Client.GetChange(ctx.Context, changeID, gerrit.QueryChangesOpt{
				Fields: []string{"CURRENT_REVISION"},
			})
			if err != nil {
				return err
			}

			branchName := fmt.Sprintf("gerrit/%d", change.ChangeNumber)
			branchRefName := plumbing.NewBranchReferenceName(branchName)

			// fetch the current changeset based on the patch number
			revision := change.Revisions[change.CurrentRevision]

			err = repo.Fetch(&git.FetchOptions{
				RemoteName: "gerrit",
				RefSpecs: []config.RefSpec{
					config.RefSpec(revision.Ref + ":" + string(branchRefName)),
				},
				Force: true,
			})
			if err != nil {
				return fmt.Errorf("failed to fetch updates for branch: %v", err)
			}

			err = wt.Checkout(&git.CheckoutOptions{
				Branch: branchRefName,
			})
			if err != nil {
				return fmt.Errorf("failed to checkout branch: %v", err)
			}

			return nil
		},
	}
)
