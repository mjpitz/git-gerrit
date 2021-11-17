package show

import (
	"strconv"
	"text/template"
	"time"

	"github.com/mjpitz/git-gerrit/internal/common"
	"github.com/mjpitz/myago/flagset"
	"github.com/urfave/cli/v2"
	"golang.org/x/build/gerrit"
)

const report = `
REVIEW:  {{ baseURL }}/c/{{ .Project }}/+/{{ .ChangeNumber }}
SUBJECT: {{ .Subject }}
OWNER:   {{ .Owner.Name }} <{{ .Owner.Email }}>

PROJECT: {{ .Project }}
BRANCH:  {{ .Branch }}

REVIEWERS:
{{- range $reviewerType, $reviewers := .Reviewers }}
{{- range $reviewer := $reviewers }}
- {{ $reviewer.Name }} <{{ $reviewer.Email }}>
{{- end }}
{{- end }}

`

type Config struct{}

var (
	showConfig = &Config{}

	Command = &cli.Command{
		Name:  "show",
		Flags: flagset.Extract(showConfig),
		Action: func(ctx *cli.Context) error {
			gerritAPI := common.GerritAPI(ctx.Context)

			accountInfo := make(map[int64]gerrit.AccountInfo)
			getAccountInfo := func(accountID int64) (*gerrit.AccountInfo, error) {
				info, ok := accountInfo[accountID]
				if !ok {
					info, err := gerritAPI.Client.GetAccountInfo(ctx.Context, strconv.FormatInt(accountID, 10))
					if err != nil {
						return nil, err
					}
					accountInfo[accountID] = info
				}

				return &info, nil
			}

			changeID := ctx.Args().Get(0)
			info, err := gerritAPI.Client.GetChangeDetail(ctx.Context, changeID)
			if err != nil {
				return err
			}

			template, err := template.New("report").
				Funcs(map[string]interface{}{
					"baseURL": func() string {
						return gerritAPI.BaseURL
					},
					"getAccountInfo": getAccountInfo,
					"formatTimeStamp": func(t gerrit.TimeStamp) string {
						return t.Time().Format(time.RFC1123Z)
					},
				}).
				Parse(report)

			if err != nil {
				return err
			}

			return template.Execute(ctx.App.Writer, info)
		},
		HideHelpCommand: true,
	}
)
