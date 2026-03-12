package run

import (
	"context"
	"fmt"
	"log"
	"time"

	libconfig "github.com/ekristen/libnuke/pkg/config"
	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	libscanner "github.com/ekristen/libnuke/pkg/scanner"
	"github.com/ekristen/libnuke/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/cgroschupp/stackit-nuke/pkg/commands/global"
	"github.com/cgroschupp/stackit-nuke/pkg/common"
	"github.com/cgroschupp/stackit-nuke/pkg/config"
	"github.com/cgroschupp/stackit-nuke/pkg/nuke"
	"github.com/cgroschupp/stackit-nuke/pkg/stackit"
	stackitclient "github.com/cgroschupp/stackit-nuke/pkg/stackit"
)

type log2LogrusWriter struct{ entry *logrus.Entry }

func (w *log2LogrusWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	w.entry.Trace(string(b))
	return n, nil
}

func execute(ctx context.Context, c *cli.Command) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	log.SetOutput(&log2LogrusWriter{entry: logrus.WithField("source", "standard-logger")})

	params := &libnuke.Parameters{
		Force:              c.Bool("force"),
		ForceSleep:         int(c.Int("force-sleep")),
		Quiet:              c.Bool("quiet"),
		NoDryRun:           c.Bool("no-dry-run"),
		Includes:           c.StringSlice("include"),
		Excludes:           c.StringSlice("exclude"),
		WaitOnDependencies: !c.Bool("no-wait-on-dependencies"),
	}

	parsedConfig, err := config.New(libconfig.Options{
		Path:         c.String("config"),
		Deprecations: registry.GetDeprecatedResourceTypeMapping(),
	})
	if err != nil {
		return err
	}

	projectID := c.String("project-id")
	if projectID == "" {
		return fmt.Errorf("--project-id is required")
	}

	region := c.String("region")

	client, err := stackitclient.NewClient(region)
	if err != nil {
		return err
	}

	filters, err := parsedConfig.Filters(projectID)
	if err != nil {
		return fmt.Errorf("parse filter: %w", err)
	}

	n := libnuke.New(params, filters, parsedConfig.Settings)
	n.SetRunSleep(5 * time.Second)
	n.SetLogger(logrus.WithField("component", "nuke"))
	n.RegisterVersion(fmt.Sprintf("> %s", common.AppVersion))

	p := &stackit.Prompt{Parameters: params, ProjectID: projectID}
	n.RegisterPrompt(p.Prompt)
	resourceTypes := types.ResolveResourceTypes(
		registry.GetNamesForScope(nuke.Account),
		[]types.Collection{n.Parameters.Includes, parsedConfig.ResourceTypes.GetIncludes()},
		[]types.Collection{n.Parameters.Excludes, parsedConfig.ResourceTypes.Excludes},
		nil,
		nil,
	)

	scanner, err := libscanner.New(&libscanner.Config{
		ResourceTypes: resourceTypes,
		Owner:         projectID,
		Opts: &nuke.ListerOpts{
			Client:    client,
			ProjectID: projectID,
			Region:    region,
			Logger:    logrus.WithField("project_id", projectID),
		},
	})
	if err != nil {
		return err
	}

	if err := n.RegisterScanner(nuke.Account, scanner); err != nil {
		return err
	}

	return n.Run(ctx)
}

func init() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "config", Usage: "path to config file", Sources: cli.Files("config.yaml")},
		&cli.StringFlag{Name: "project-id", Usage: "STACKIT project UUID", Required: true},
		&cli.StringFlag{Name: "region", Usage: "STACKIT region", Required: true},
		&cli.StringSliceFlag{Name: "include", Usage: "only include a specific resource type"},
		&cli.StringSliceFlag{Name: "exclude", Usage: "exclude a specific resource type"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "hide filtered messages"},
		&cli.BoolFlag{Name: "no-dry-run", Usage: "actually remove resources after discovery"},
		&cli.BoolFlag{Name: "no-prompt", Usage: "disable prompting for verification", Aliases: []string{"force"}},
		&cli.BoolFlag{Name: "no-wait-on-dependencies", Usage: "disable dependency waiting"},
		&cli.IntFlag{Name: "prompt-delay", Usage: "seconds to delay after prompt before running", Value: 10, Aliases: []string{"force-sleep"}},
	}

	cmd := &cli.Command{
		Name:    "run",
		Aliases: []string{"nuke"},
		Usage:   "run nuke against a STACKIT project",
		Flags:   append(flags, global.Flags()...),
		Before:  global.Before,
		Action:  execute,
	}

	common.RegisterCommand(cmd)
}
