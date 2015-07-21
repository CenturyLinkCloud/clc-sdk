package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/sdk/alert"
	"github.com/mikebeyer/clc-sdk/sdk/clc"
)

// Commands exports the cli commands for the status package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "alert",
		Aliases: []string{"a"},
		Usage:   "alert api",
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
		Aliases: []string{"a"},
		Usage:   "get alert policy",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "retrieve all policies"},
			cli.StringFlag{Name: "alias, a", Usage: "policy id"},
		},
		Action: func(c *cli.Context) {
			if c.Bool("all") || c.Args().First() == "" {
				policies, err := client.Alert.GetAll()
				if err != nil {
					fmt.Printf("failed to get %s", c.Args().First())
					return
				}
				b, err := json.MarshalIndent(policies, "", "  ")
				if err != nil {
					fmt.Printf("%s", err)
				}
				fmt.Printf("%s\n", b)
				return
			}
			policy, err := client.Alert.Get(c.Args().First())
			if err != nil {
				fmt.Printf("failed to get %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(policy, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func create(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create alert policy",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "policy name [required]"},
			cli.StringSliceFlag{Name: "email, e", Usage: "provide email address for alert"},
			cli.StringSliceFlag{Name: "cpu, c", Usage: "provide threshold / duration in minutes (ex. --cpu 90/10)"},
			cli.StringSliceFlag{Name: "disk, d", Usage: "provide threshold / duration in minutes (ex. --disk 90/10)"},
			cli.StringSliceFlag{Name: "memory, m", Usage: "provide threshold / duration in minutes (ex. --memory 90/10)"},
		},
		Before: func(c *cli.Context) error {
			if c.String("name") == "" {
				fmt.Printf("missing flags, --help for usage")
				return errors.New("")
			}

			return nil
		},
		Action: func(c *cli.Context) {
			a := alert.Alert{
				Name: c.String("name"),
			}
			if len(c.StringSlice("email")) > 0 {
				action := alert.Action{
					Action: "email",
					Setting: alert.Setting{
						Recipients: c.StringSlice("email"),
					},
				}
				a.Actions = []alert.Action{action}
			}

			cpus, err := parseTrigger("cpu", c.StringSlice("cpu"))
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			a.Triggers = append(a.Triggers, cpus...)

			disks, err := parseTrigger("disk", c.StringSlice("disk"))
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			a.Triggers = append(a.Triggers, disks...)

			mems, err := parseTrigger("memory", c.StringSlice("memory"))
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			a.Triggers = append(a.Triggers, mems...)
			policy, err := client.Alert.Create(a)
			if err != nil {
				log.Printf("err: %s\n", err)
				fmt.Printf("failed to create policy")
				return
			}

			b, err := json.MarshalIndent(policy, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func parseTrigger(name string, metrics []string) ([]alert.Trigger, error) {
	triggers := make([]alert.Trigger, 0)
	for _, v := range metrics {
		split := strings.Split(v, "/")
		threshold, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return triggers, fmt.Errorf("failed to parse cpu flag %s", v)
		}

		dur, err := strconv.Atoi(split[1])
		if err != nil {
			return triggers, fmt.Errorf("failed to parse cpu flag %s", v)
		}

		trigger := alert.Trigger{
			Metric:    name,
			Duration:  convertMinutesToTimeFormat(dur),
			Threshold: threshold,
		}
		triggers = append(triggers, trigger)
	}
	return triggers, nil
}

func convertMinutesToTimeFormat(amount int) string {
	hrs := amount / 60
	mins := amount - hrs*60
	return fmt.Sprintf("%s:%s:00", twoDigitTime(hrs), twoDigitTime(mins))
}

func twoDigitTime(amount int) string {
	if amount < 0 {
		return "00"
	}
	if amount < 10 {
		return fmt.Sprintf("0%v", amount)
	} else {
		return fmt.Sprintf("%v", amount)
	}
}

func delete(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete alert policy",
		Before: func(c *cli.Context) error {
			if c.Args().First() == "" {
				fmt.Println("usage: delete [policy]")
				return errors.New("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			err := client.Alert.Delete(c.Args().First())
			if err != nil {
				fmt.Printf("failed to get %s", c.Args().First())
				os.Exit(1)
			}
			fmt.Printf("deleted policy %s\n", c.Args().First())
		},
	}
}
