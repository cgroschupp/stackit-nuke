package main

import (
	"context"
	"net/mail"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/cgroschupp/stackit-nuke/pkg/commands/list"
	_ "github.com/cgroschupp/stackit-nuke/pkg/commands/run"
	"github.com/cgroschupp/stackit-nuke/pkg/common"
	_ "github.com/cgroschupp/stackit-nuke/resources"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.Command{}
	app.Authors = []any{&mail.Address{Name: "Christian", Address: "christian@groschupp.org"}}
	app.Name = path.Base(os.Args[0])
	app.Usage = "remove supported resources from a STACKIT project"
	app.Version = common.AppVersion.Summary
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(ctx context.Context, c *cli.Command, command string) {
		logrus.Fatalf("Command %s not found.", command)
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		logrus.Fatalf("unable to run: %s", err)
	}
}
