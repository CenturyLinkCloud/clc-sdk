package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
)

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
				fmt.Errorf("failed to get status of %s", c.Args().First())
				log.Printf("%s", err)
				os.Exit(1)
			}
			b, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
