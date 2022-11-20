package main

import (
	"filesystem-api/api"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	fsAPICfg = api.Config{}
	fsAPI    = api.Handler{}
)

func main() {
	app := cli.App{
		Name:        "Filesystem API",
		Description: "api used to interact with a remote filesystem",
	}
	app.Flags = append(app.Flags, fsAPICfg.Flags()...)
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	fsAPI.Initialize(&fsAPICfg)
	return fsAPI.Start()
}
