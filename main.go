package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

var version = time.Now().String()

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Name = "gdom"
	app.Usage = `Google G Suite command line tool

	see https://github.com/emicklei/gdom for documentation.
`
	// override -v
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose logging",
		},
	}
	format := cli.StringFlag{
		Name:  "format",
		Usage: "-format JSON",
	}
	app.Commands = []cli.Command{
		{
			Name:  "user",
			Usage: "Retrieving information related to user accounts",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "Show list of all users",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit",
							Usage: "-limit 10",
						},
						format,
					},
					Action: func(c *cli.Context) error {
						return cmdUserList(c)
					},
					ArgsUsage: `user list`,
				},
				{
					Name:  "membership",
					Usage: "Show list of groups for which the user has a membership",
					Action: func(c *cli.Context) error {
						return cmdUserMembershipList(c)
					},
					ArgsUsage: `user membership john.doe@company.com`,
				},
				{
					Name:  "info",
					Usage: "Show user details",
					Action: func(c *cli.Context) error {
						return cmdUserInfo(c)
					},
					Flags:     []cli.Flag{format},
					ArgsUsage: `user info john.doe@company.com`,
				},
			},
		},
		{
			Name:  "group",
			Usage: "Retrieving information related to groups",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "Show list of all groups",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit",
							Usage: "-limit 10",
						},
						format,
					},
					Action: func(c *cli.Context) error {
						return cmdGroupList(c)
					},
					ArgsUsage: `group list`,
				},
				{
					Name:  "members",
					Usage: "Show members of a group",
					Action: func(c *cli.Context) error {
						return cmdGroupMembers(c)
					},
					Flags:     []cli.Flag{format},
					ArgsUsage: `group members all@company.com`,
				},
				{
					Name:  "info",
					Usage: "Show group details",
					Action: func(c *cli.Context) error {
						return cmdGroupInfo(c)
					},
					Flags:     []cli.Flag{format},
					ArgsUsage: `user info all@company.com`,
				},
			},
		},
		{
			Name:  "reset",
			Usage: "Forget about the cached credentials and scopes",
			Action: func(c *cli.Context) error {
				if c.GlobalBool("v") {
					fmt.Println("[gdom] delete $HOME/gdom-token.json (if present)")
				}
				return os.Remove(filepath.Join(os.Getenv("HOME"), "gdom-token.json"))
			},
		},
	}
	return app
}
