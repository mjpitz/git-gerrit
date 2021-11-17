package log

import (
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/mjpitz/git-gerrit/internal/common"
	"github.com/mjpitz/myago/flagset"
	"github.com/urfave/cli/v2"
	"golang.org/x/build/gerrit"
)

const log = `
{{- range $message := .Messages }}
---
Author: {{ $message.Author.Name }} <{{ $message.Author.Email }}>
Date:   {{ formatTimeStamp $message.Time }}

{{ $message.Message }}
{{ end }}
`

type Config struct {
}

var (
	config = &Config{}

	Command = &cli.Command{
		Name:      "log",
		Usage:     "Output updates about a changeset",
		UsageText: "log [options] <changeID>",
		Flags:     flagset.Extract(config),
		Action: func(ctx *cli.Context) error {
			gerritAPI := common.GerritAPI(ctx.Context)

			// TODO: infer from current "branch"
			changeID := ctx.Args().Get(0)

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

			changeInfo, err := gerritAPI.Client.GetChangeDetail(ctx.Context, changeID)
			if err != nil {
				return err
			}

			sort.Slice(changeInfo.Messages, func(i, j int) bool {
				a := changeInfo.Messages[i].Time.Time()
				b := changeInfo.Messages[j].Time.Time()
				return b.Before(a)
			})

			for _, message := range changeInfo.Messages {
				message.Author, _ = getAccountInfo(message.Author.NumericID)
			}

			template, err := template.New("log").
				Funcs(map[string]interface{}{
					"baseURL": func() string {
						return gerritAPI.BaseURL
					},
					"getAccountInfo": getAccountInfo,
					"formatTimeStamp": func(t gerrit.TimeStamp) string {
						return t.Time().Local().Format(time.RFC1123Z)
					},
				}).
				Parse(log)

			if err != nil {
				return err
			}

			return template.Execute(ctx.App.Writer, changeInfo)
		},
		HideHelpCommand: true,
	}
)
