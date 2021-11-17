package checkout

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/go-git/go-git/v5"
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

			changeID := ctx.Args().Get(0)

			change, err := gerritAPI.Client.GetChange(ctx.Context, changeID, gerrit.QueryChangesOpt{
				Fields: []string{"CURRENT_REVISION"},
			})

			if err != nil {
				return err
			}

			wt, err := repo.Worktree()
			if err != nil {
				return fmt.Errorf("failed to obtain worktree: %v", err)
			}

			revision := plumbing.NewHash(change.CurrentRevision)
			err = wt.Checkout(&git.CheckoutOptions{
				Hash: revision,
			})

			if err != nil {
				return fmt.Errorf("failed to checkout current revision: %v", err)
			}

			branchName := fmt.Sprintf("gerrit/%d", change.ChangeNumber)
			err = repo.DeleteBranch(branchName)
			switch {
			case errors.Is(err, git.ErrBranchNotFound):
			case err != nil:
				return fmt.Errorf("failed to delete last branch: %v", err)
			}

			cmd := fmt.Sprintf("git branch -D %s; git checkout -b %s", branchName, branchName)
			err = exec.Command("bash", "-c", cmd).Run()
			if err != nil {
				return fmt.Errorf("failed to checkout new branch: %v", err)
			}

			cmd = fmt.Sprintf("git rebase %s", change.Branch)
			err = exec.Command("bash", "-c", cmd).Run()
			if err != nil {
				return fmt.Errorf("failed to rebase new branch: %v", err)
			}

			return nil
		},
	}
)
