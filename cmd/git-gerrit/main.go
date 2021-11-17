package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mjpitz/git-gerrit/internal/commands/changes"
	"github.com/mjpitz/git-gerrit/internal/commands/log"
	"github.com/mjpitz/git-gerrit/internal/commands/show"
	"github.com/mjpitz/git-gerrit/internal/common"
	"github.com/urfave/cli/v2"

	"github.com/mjpitz/git-gerrit/internal/commands"
	"github.com/mjpitz/myago/flagset"
	"github.com/mjpitz/myago/lifecycle"
	"github.com/mjpitz/myago/zaputil"
)

var version = ""
var commit = ""
var date = time.Now().Format(time.RFC3339)

type GlobalConfig struct {
	Log zaputil.Config `json:"log"`
}

func main() {
	compiled, _ := time.Parse(time.RFC3339, date)

	cfg := &GlobalConfig{
		Log: zaputil.DefaultConfig(),
	}

	app := &cli.App{
		Name:      "git-gerrit",
		Usage:     "Gerrit <3 Git",
		UsageText: "git-gerrit [options] <command>",
		Version:   fmt.Sprintf("%s (%s)", version, commit),
		Flags:     flagset.Extract(cfg),
		Commands: []*cli.Command{
			changes.Command,
			log.Command,
			show.Command,
			commands.Version,
		},
		Before: func(ctx *cli.Context) (err error) {
			ctx.Context = zaputil.Setup(ctx.Context, cfg.Log)
			ctx.Context = lifecycle.Setup(ctx.Context)
			ctx.Context = common.SetupTableWriter(ctx.Context, ctx.App.Writer)

			ctx.Context, err = common.SetupGitRepository(ctx.Context)
			if err != nil {
				return err
			}

			ctx.Context, err = common.SetupProject(ctx.Context)
			if err != nil {
				return err
			}

			ctx.Context, err = common.SetupGerritAPI(ctx.Context)
			return err
		},
		After: func(ctx *cli.Context) error {
			lifecycle.Resolve(ctx.Context)
			common.TableWriter(ctx.Context).Render()

			return nil
		},
		Compiled:             compiled,
		Copyright:            fmt.Sprintf("Copyright %d The git-gerrit Authors - All Rights Reserved\n", compiled.Year()),
		HideVersion:          true,
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		BashComplete:         cli.DefaultAppComplete,
		Metadata: map[string]interface{}{
			"arch":       runtime.GOARCH,
			"compiled":   date,
			"go_version": strings.TrimPrefix(runtime.Version(), "go"),
			"os":         runtime.GOOS,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
