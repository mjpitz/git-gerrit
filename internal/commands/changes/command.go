package changes

import (
	"sort"
	"time"

	"github.com/mjpitz/git-gerrit/internal/common"
	"github.com/mjpitz/myago/flagset"
	"github.com/mjpitz/myago/zaputil"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type Config struct {
	AllProjects bool   `json:"all_projects" alias:"A" usage:""`
	Query       string `json:"query"        alias:"q" usage:""`
}

type Change struct {
	Updated  time.Time
	ChangeID string `json:"change_id"`
	Subject  string
}

var (
	config = &Config{}

	Command = &cli.Command{
		Name:  "changes",
		Flags: flagset.Extract(config),
		Action: func(ctx *cli.Context) error {
			gerritAPI := common.GerritAPI(ctx.Context)
			projectID := common.Project(ctx.Context)
			writer := common.TableWriter(ctx.Context)

			if config.Query == "" {
				config.Query = "status:open -is:wip"
				if !config.AllProjects {
					config.Query = config.Query + " project:" + string(projectID)
				}
			}

			log := zaputil.Extract(ctx.Context)

			log.Debug("query changes", zap.String("q", config.Query))
			changes, err := gerritAPI.Client.QueryChanges(ctx.Context, config.Query)
			if err != nil {
				return err
			}

			sort.Slice(changes, func(i, j int) bool {
				return changes[j].Updated.Time().Before(changes[i].Updated.Time())
			})

			for _, change := range changes {
				common.WriteRow(writer, &Change{
					Updated:  time.Time(change.Updated).Local(),
					ChangeID: change.ChangeID,
					Subject:  change.Subject,
				})
			}

			return nil
		},
		HideHelpCommand: true,
	}
)
