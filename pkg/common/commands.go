package common

import "github.com/urfave/cli/v3"

var commands []*cli.Command

func RegisterCommand(c *cli.Command) {
	commands = append(commands, c)
}

func GetCommands() []*cli.Command {
	return commands
}
