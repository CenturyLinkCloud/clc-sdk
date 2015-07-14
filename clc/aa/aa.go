package aa

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/sdk/clc"
)

// Commands exports the cli commands for the status package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:        "anti-alias",
		Aliases:     []string{"aa"},
		Usage:       "anti-alias api",
		Subcommands: []cli.Command{get(client)},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get aa policy",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "retrieve all policies"},
			cli.StringFlag{Name: "alias, a", Usage: "account alias"},
		},
		Action: func(c *cli.Context) {
			if c.Bool("all") || c.Args().First() == "" {
				policies, err := client.AA.GetAll()
				if err != nil {
					fmt.Printf("unable to retrieve aa policies")
					return
				}
				b, err := json.MarshalIndent(policies, "", "  ")
				if err != nil {
					fmt.Printf("%s", err)
					return
				}
				fmt.Printf("%s\n", b)
			}

			policy, err := client.AA.Get(c.Args().First())
			if err != nil {
				fmt.Printf("unable to retrieve aa policy: [%s]", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(policy, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}
