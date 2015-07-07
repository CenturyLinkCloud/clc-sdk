package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
)

func main() {
	client := clc.New(clc.EnvConfig())

	app := cli.NewApp()
	app.Name = "clc"
	app.Usage = "v2 api"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "interact with server api",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get [name] server details",
					Before: func(c *cli.Context) error {
						if c.Args().First() == "" {
							fmt.Println("usage: get [name]")
							return errors.New("")
						}
						return nil
					},
					Action: func(c *cli.Context) {
						server, err := client.Server.Get(c.Args().First())
						if err == nil {
							b, err := json.MarshalIndent(server, "", "  ")
							if err == nil {
								fmt.Printf("%s\n", b)
							}
						}
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
