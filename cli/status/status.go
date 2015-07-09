package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
				log.Fatalf("failed to get status of %s", c.Args().First())
			}
			b, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				log.Fatalf("%s", err)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
