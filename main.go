package main

import (
	"log"
	"github.com/urfave/cli/v2"
	"os"
)
// cd go-docker
// $GOROOT/bin/go build .
// ./go-docker run -it /bin/sh
// ./go-docker run -it /bin/ls



func main() {
	app := &cli.App {
		Name: "go-docker",
		Usage: "go-docker run -it xxx",
	}
	app.Commands = []*cli.Command {
		&initCommand,
		&runCommand,
	}
	app.Before = func(c *cli.Context) error {
		// log.SetFormatter(&log.JSONFormatter)
		log.SetOutput(os.Stdout)

		return nil
	}

	if err := app.Run(os.Args); err!= nil {
		log.Fatal(err)
	}

}


