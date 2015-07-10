package status

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
)

// Commands exports the cli commands for the status package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:        "status",
		Usage:       "status api",
		Subcommands: []cli.Command{get(client)},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get status",
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				fmt.Println("usage: get [id]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			status, err := client.Status.Get(c.Args().First(), nil)
			if err != nil {
				fmt.Printf("failed to get status of %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}
