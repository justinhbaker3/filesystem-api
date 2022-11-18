package main

import (
	"filesystem-api/api"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:        "Filesystem API",
		Description: "api used to interact with a remote filesystem",
	}

	h := api.NewHandler()
	app.Flags = append(app.Flags, h.Flags()...)
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {

}
