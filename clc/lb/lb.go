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
			getPool(client),
			createPool(client),
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
		Before: func(c *cli.Context) error {
			if c.String("location") == "" {
				fmt.Printf("location flag is required.\n")
				return fmt.Errorf("")
			}
			return nil
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
				fmt.Printf("missing required flags to create load balancer. [use --help to show required flags]\n")
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

func getPool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get-pool",
		Aliases: []string{"gp"},
		Usage:   "get load balancer pool details",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "list all load balancers for location"},
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool", Usage: "load balancer pool id"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" || c.String("id") == "" {
				fmt.Printf("missing required flags to get pool. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			if c.Bool("all") || c.String("pool") == "" {
				resp, err := client.LB.GetAllPools(c.String("location"), c.String("id"))
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
			resp, err := client.LB.GetPool(c.String("location"), c.String("id"), c.String("pool"))
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

func createPool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create-pool",
		Aliases: []string{"cp"},
		Usage:   "create load balancer pool",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.IntFlag{Name: "port", Usage: "pool port [required]"},
			cli.BoolFlag{Name: "sticky", Usage: "use stick persistence"},
			cli.BoolFlag{Name: "standard", Usage: "use standard persistence [default]"},
			cli.BoolFlag{Name: "least-connection, lc", Usage: "use least-connection load balacing"},
			cli.BoolFlag{Name: "round-robin, rr", Usage: "use round-robin load balacing [default]"},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("sticky") && c.Bool("standard") {
				fmt.Println("only one of sticky and standard can be selected")
				return fmt.Errorf("")
			}

			if c.Bool("least-connection") && c.Bool("round-robin") {
				fmt.Println("only one of least-connection and round-robin can be selected")
				return fmt.Errorf("")
			}

			if c.String("id") == "" || c.String("location") == "" || c.Int("port") == 0 {
				fmt.Println("missing required flags, --help for more details")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			pool := lb.Pool{Port: c.Int("port")}
			if c.Bool("sticky") {
				pool.Persistence = lb.Sticky
			} else {
				pool.Persistence = lb.Standard
			}

			if c.Bool("least-connection") {
				pool.Method = lb.LeastConn
			} else {
				pool.Method = lb.RoundRobin
			}

			resp, err := client.LB.CreatePool(c.String("location"), c.String("id"), pool)
			if err != nil {
				fmt.Printf("failed to create load balancer pool for [%s] in %s\n", c.String("id"), c.String("location"))
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
