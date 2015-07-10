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

// Commands exports the cli commands for the server package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"s"},
		Usage:       "server api",
		Subcommands: []cli.Command{get(client), create(client), delete(client)},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get server details",
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
				log.Fatalf("failed to get %s", c.Args().First())
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
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
		Usage:   "create server",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringFlag{Name: "cpu, c", Usage: "server cpus (1 - 16) [required]"},
			cli.StringFlag{Name: "memory, m", Usage: "server memory in gbs (1 - 128) [required]"},
			cli.StringFlag{Name: "group, g", Usage: "parent group id [required]"},
			cli.StringFlag{Name: "source, s", Usage: "source server id (template or existing server) [required]"},
			cli.StringFlag{Name: "type, t", Usage: "standard or hyperscale [required]"},
			cli.StringFlag{Name: "password, p", Usage: "server password"},
			cli.StringFlag{Name: "description, d", Usage: "server description"},
			cli.StringFlag{Name: "ip", Usage: "id address"},
			cli.BoolFlag{Name: "managed", Usage: "make server managed"},
			cli.StringFlag{Name: "primaryDNS", Usage: "primary dns"},
			cli.StringFlag{Name: "secondaryDNS", Usage: "secondary dns"},
			cli.StringFlag{Name: "network", Usage: "network id"},
			cli.StringFlag{Name: "storage", Usage: "standard or premium"},
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
			server.Password = c.String("password")
			server.Description = c.String("description")
			server.IPaddress = c.String("ip")
			server.IsManagedOS = c.Bool("managed")
			server.PrimaryDNS = c.String("primaryDNS")
			server.SecondaryDNS = c.String("secondaryDNS")
			server.NetworkID = c.String("network")
			server.Storagetype = c.String("storage")

			if !server.Valid() {
				fmt.Println("missing required flags to create server. [use --help to show required flags]")
				return
			}

			resp, err := client.Server.Create(server, nil)
			if err != nil {
				fmt.Printf("failed to create %s", server.Name)
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
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
		Usage:   "delete server",
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				fmt.Println("usage: delete [server]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			server, err := client.Server.Delete(c.Args().First())
			if err != nil {
				log.Fatalf("failed to delete %s", c.Args().First())
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
