package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/mikebeyer/clc-sdk/cli/server"
	"github.com/mikebeyer/clc-sdk/cli/status"
)

func main() {
	client := clc.New(clc.EnvConfig())

	app := cli.NewApp()
	app.Name = "clc"
	app.Usage = "clc v2 api cli"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mike Beyer",
			Email: "michael.beyer@centurylink.com",
		},
	}
	app.Commands = []cli.Command{
		test(client),
		server.Commands(client),
		status.Commands(client),
	}
	app.Run(os.Args)
}

func test(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Action: func(c *cli.Context) {
			token, err := client.Auth()
			if err != nil {
				fmt.Printf("test failed [%s]", err)
				return
			}
			fmt.Printf("success: [%s]", token)
		},
	}
}
