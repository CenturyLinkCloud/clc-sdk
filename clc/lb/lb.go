package lb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/sdk/clc"
	"github.com/mikebeyer/clc-sdk/sdk/lb"
)

func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "load-balancer",
		Aliases: []string{"lb"},
		Usage:   "load balancer api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
		},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get load balancer details",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "list all load balancers for location"},
			cli.StringFlag{Name: "id", Usage: "load balancer id"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
		},
		Action: func(c *cli.Context) {
			if c.Bool("all") || c.String("id") == "" {
				resp, err := client.LB.GetAll(c.String("location"))
				if err != nil {
					log.Fatalf("failed to get %s\n", c.Args().First())
				}
				b, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					log.Printf("%s\n", err)
					os.Exit(1)
				}
				fmt.Printf("%s\n", b)
				return
			}
			resp, err := client.LB.Get(c.String("location"), c.String("id"))
			if err != nil {
				log.Fatalf("failed to get %s\n", c.Args().First())
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				log.Printf("%s\n", err)
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
		Usage:   "create shared load balancer",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "load balancer name [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "description, d", Usage: "load balancer description"},
		},
		Action: func(c *cli.Context) {
			name := c.String("name")
			loc := c.String("location")
			if name == "" || loc == "" {
				fmt.Printf("missing required flags to load balancer policy. [use --help to show required flags]\n")
				return
			}

			lb := lb.LoadBalancer{Name: name, Description: c.String("description")}
			resp, err := client.LB.Create(loc, lb)
			if err != nil {
				fmt.Printf("failed to create load balancer [%s] in %s\n", name, loc)
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}
