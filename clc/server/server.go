package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/sdk/clc"
	"github.com/mikebeyer/clc-sdk/sdk/server"
)

// Commands exports the cli commands for the server package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "server api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
			delete(client),
			publicIP(client),
			archive(client),
			restore(client),
		},
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
				fmt.Printf("failed to get %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(server, "", "  ")
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
			server := server.Server{
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

			resp, err := client.Server.Create(server)
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
				fmt.Printf("failed to delete %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func archive(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "archive",
		Aliases: []string{"a"},
		Usage:   "archive server",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "name, n", Usage: "name of servers to archive"},
		},
		Before: func(c *cli.Context) error {
			if len(c.StringSlice("name")) == 0 {
				fmt.Println("usage: -n [server] -n [server]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			resp, err := client.Server.Archive(c.StringSlice("name")...)
			if err != nil {
				fmt.Printf("failed to archive %s", strings.Join(c.StringSlice("name"), ", "))
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

func restore(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "restore",
		Aliases: []string{"r"},
		Usage:   "restore server",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "name of server to restore [required]"},
			cli.StringFlag{Name: "group, g", Usage: "group for server to restore to [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("name") == "" || c.String("group") == "" {
				fmt.Println("missing required flags [--help for additional information]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			resp, err := client.Server.Restore(c.String("name"), c.String("group"))
			if err != nil {
				fmt.Printf("failed to restore %s", c.String("name"))
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

func publicIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "public-ip",
		Aliases: []string{"ip"},
		Usage:   "manage public ips",
		Subcommands: []cli.Command{
			getIP(client),
			createIP(client),
			deleteIP(client),
		},
	}
}

func getIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get public ip",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringFlag{Name: "ip", Usage: "ip [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("name") == "" || c.String("ip") == "" {
				fmt.Println("usage: missing required flags [--help for additional information]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			resp, err := client.Server.GetPublicIP(c.String("name"), c.String("ip"))
			if err != nil {
				fmt.Printf("err %s\n", err)
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

func createIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add public ip to server",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringSliceFlag{Name: "tcp", Usage: "provide a port [8080] or a port range [8080:8082]"},
			cli.StringSliceFlag{Name: "udp", Usage: "provide a port [8080] or a port range [8080:8082]"},
			cli.StringSliceFlag{Name: "restriction, r", Usage: "provide an ip subnet to restrict to access the public ip [ex. 10.0.0.1/24 (must be cidr notation)]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("name") == "" {
				fmt.Println("usage: --name flag required")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			ports := make([]server.Port, 0)
			tcps, err := parsePort("tcp", c.StringSlice("tcp"))
			if err != nil {
				fmt.Println(err.Error())
			}
			ports = append(ports, tcps...)
			udps, err := parsePort("udp", c.StringSlice("udp"))
			if err != nil {
				fmt.Println(err.Error())
			}
			ports = append(ports, udps...)
			restrictions := make([]server.SourceRestriction, 0)
			for _, v := range c.StringSlice("restriction") {
				restrictions = append(restrictions, server.SourceRestriction{CIDR: v})
			}

			ip := server.PublicIP{Ports: ports}
			resp, err := client.Server.AddPublicIP(c.String("name"), ip)
			if err != nil {
				fmt.Printf("err %s\n", err)
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

func parsePort(protocol string, list []string) ([]server.Port, error) {
	ports := make([]server.Port, 0)
	for _, v := range list {
		r := strings.Split(v, ":")
		port, err := strconv.Atoi(r[0])
		if err != nil {
			return ports, fmt.Errorf("invalid port provided %s", r[0])
		}
		if len(r) > 1 {
			to, err := strconv.Atoi(r[1])
			if err != nil {
				return ports, fmt.Errorf("invalid port provided %s", r[0])
			}
			ports = append(ports, server.Port{Protocol: protocol, Port: port, PortTo: to})
		} else {
			ports = append(ports, server.Port{Protocol: protocol, Port: port})
		}
	}
	return ports, nil
}

func deleteIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete public ip",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringFlag{Name: "ip", Usage: "ip [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("name") == "" || c.String("ip") == "" {
				fmt.Println("usage: missing required flags [--help for additional information]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			resp, err := client.Server.DeletePublicIP(c.String("name"), c.String("ip"))
			if err != nil {
				fmt.Printf("err %s\n", err)
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
