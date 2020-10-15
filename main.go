package main

import (
	"log"
	"os"

	"github.com/moonguard-org/moonguard/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "moonguard",
		Usage: "gRPC tooling",
		Commands: []*cli.Command{
			commands.GetGenCommand(),
			commands.GetInitializeCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
