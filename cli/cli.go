package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/mikebeyer/clc-sdk/cli/server"
)

func main() {
	client := clc.New(clc.EnvConfig())

	app := cli.NewApp()
	app.Name = "clc"
	app.Usage = "v2 api"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{server.Commands(client)}
	app.Run(os.Args)
}
