package server

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
		Name:        "server",
		Aliases:     []string{"s"},
		Usage:       "interact with server api",
		Subcommands: []cli.Command{get(client), create(client)},
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
				log.Printf("%s", err)
				os.Exit(1)
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				fmt.Errorf("failed to get %s", c.Args().First())
				log.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func create(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create ",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "server name",
			},
			cli.StringFlag{
				Name:  "cpu, c",
				Usage: "server cpus (1 - 16)",
			},
			cli.StringFlag{
				Name:  "memory, m",
				Usage: "server memory in gbs (1 - 128)",
			},
			cli.StringFlag{
				Name:  "group, g",
				Usage: "parent group id",
			},
			cli.StringFlag{
				Name:  "source, s",
				Usage: "source server id (template or existing server)",
			},
			cli.StringFlag{
				Name:  "type, t",
				Usage: "standard or hyperscale",
			},
		},
		Action: func(c *cli.Context) {
			server := clc.Server{
				Name:           c.String("name"),
				CPU:            c.Int("cpu"),
				MemoryGB:       c.Int("memory"),
				GroupID:        c.String("group"),
				SourceServerID: c.String("source"),
				Type:           c.String("type"),
			}
			resp, err := client.Server.Create(server, nil)
			if err != nil {
				fmt.Errorf("failed to create %s", server.Name)
				log.Printf("%s", err)
				os.Exit(1)
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Errorf("failed to create %s", server.Name)
				log.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
