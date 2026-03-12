package main

import (
	"fmt"

	"github.com/wendisx/puzzle/internal/command"
	"github.com/wendisx/puzzle/pkg/cli"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
	"github.com/wendisx/puzzle/pkg/palette"
)

func LoadDict(dictkey string) {
	if config.HasDict(config.DictKey(dictkey)) {
		return
	}
	dict := config.NewDataDict[any](dictkey)
	config.PutDict(dict.Name(), dict)
	clog.Info(fmt.Sprintf("load dict(%s) successfully", palette.SkyBlue(dictkey)))
}

func main() {
	clog.DefaultLevel(clog.WARN)
	_ = config.LoadConfig("./demo/dev.yaml")
	LoadDict(config.DICTKEY_COMMAND)
	cli.Execute(
		command.MountBuiltinVersion,
		command.MountBuiltinInit,
		command.MountBuiltinNew,
	)
}
