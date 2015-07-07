package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
)

func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"s"},
		Usage:       "interact with server api",
		Subcommands: []cli.Command{get(client)},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get [server] server details",
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				fmt.Println("usage: get [server]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			server, err := client.Server.Get(c.Args().First())
			if err != nil {
				fmt.Errorf("failed to get %s", c.Args().First())
				os.Exit(1)
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				fmt.Errorf("failed to get %s", c.Args().First())
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
