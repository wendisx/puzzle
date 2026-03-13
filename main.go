package main

import (
	"github.com/wendisx/puzzle/internal/command"
	"github.com/wendisx/puzzle/pkg/cli"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

func main() {
	clog.DefaultLevel(clog.WARN)
	_ = config.LoadConfig("./demo/dev.yaml")
	config.LoadDict(config.DICTKEY_COMMAND)
	cli.Execute(
		command.MountBuiltinVersion,
		command.MountBuiltinInit,
		command.MountBuiltinNew,
	)
}
