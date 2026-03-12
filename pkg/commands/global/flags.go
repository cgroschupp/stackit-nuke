package global

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{Name: "debug", Usage: "enable debug logging"},
		&cli.BoolFlag{Name: "trace", Usage: "enable trace logging"},
	}
}

func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	switch {
	case c.Bool("trace"):
		logrus.SetLevel(logrus.TraceLevel)
	case c.Bool("debug"):
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	return ctx, nil
}
