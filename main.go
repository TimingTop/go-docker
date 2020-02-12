package main

import (
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	"os"
)

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


