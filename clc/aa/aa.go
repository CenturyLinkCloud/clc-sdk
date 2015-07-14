package aa

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/sdk/clc"
)

// Commands exports the cli commands for the status package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "anti-alias",
		Aliases: []string{"aa"},
		Usage:   "anti-alias api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
			delete(client),
		},
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
					fmt.Printf("unable to retrieve aa policies\n")
					return
				}
				b, err := json.MarshalIndent(policies, "", "  ")
				if err != nil {
					fmt.Printf("%s", err)
					return
				}
				fmt.Printf("%s\n", b)
				return
			}

			policy, err := client.AA.Get(c.Args().First())
			if err != nil {
				fmt.Printf("unable to retrieve aa policy: [%s]\n", c.Args().First())
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

func create(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create aa policy",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "policy name [required]"},
			cli.StringFlag{Name: "location, l", Usage: "policy location [required]"},
		},
		Action: func(c *cli.Context) {
			name := c.String("name")
			loc := c.String("location")
			if name == "" || loc == "" {
				fmt.Printf("missing required flags to create policy. [use --help to show required flags]\n")
				return
			}

			policy, err := client.AA.Create(name, loc)
			if err != nil {
				fmt.Printf("failed to create policy %s in %s", name, loc)
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

func delete(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete aa policy",
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				fmt.Println("usage: delete [id]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			err := client.AA.Delete(c.Args().First())
			if err != nil {
				fmt.Printf("unable to delete aa policy: [%s]\n", c.Args().First())
				return
			}
			fmt.Printf("deleted aa policy: %s\n", c.Args().First())
		},
	}
}
