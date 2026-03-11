package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

// test basic load cmd [passed]
func Test_load_cmd(t *testing.T) {
	// test load cmd [passed]
	config.LoadConfig("../../demo/dev.yaml")
	cli := LoadCmd("")
	clog.Info(fmt.Sprintf("%#v", cli))
	// test load cobra command [passed]
	cmdDict := config.GetDict(config.DICTKEY_COMMAND)
	verb := "client:start"
	delimiter := ":"
	cmd, ok := cmdDict.Find(strings.ReplaceAll(verb, delimiter, "")).Value().(*cobra.Command)
	if !ok {
		clog.Panic(fmt.Sprintf("from dict_key(%s) assert to type(*cobra.Command) fail", strings.ReplaceAll(verb, delimiter, "")))
	}
	clog.Info(fmt.Sprintf("%#v", cmd))
	// test relative for client and client:start [passed]
	verbP := "client"
	cmdP, ok := cmdDict.Find(strings.ReplaceAll(verbP, delimiter, "")).Value().(*cobra.Command)
	if !ok {
		clog.Panic(fmt.Sprintf("from dict_key(%s) assert to type(*cobra.Command) fail", strings.ReplaceAll(verbP, delimiter, "")))
	}
	clog.Info(fmt.Sprintf("%p", cmdP))
	// test fail [passed]
	cli = LoadCmd("./command.json")
	clog.Info(fmt.Sprintf("%#v", cli))
}

// test get flags [passed]
func Test_get_flags(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	cli := LoadCmd("")
	// test persistent flags [passed]
	delimiter := ":"
	verb := ":server:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	// test local flags [passed]
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// test diff verb with the same sub verb [passed]
	verb = ":client:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// expected [2,1,0,3] [passed]
	// not exists verb [passed]
	verb = ":log:start"
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
	// test empty verb [passed]
	verb = ""
	clog.Info(fmt.Sprintf("%#v", ParsePersistenFlags(verb, delimiter, cli)))
	clog.Info(fmt.Sprintf("%#v", ParseLocalFlags(verb, delimiter, cli)))
}

// test get cli from data dict  [passed]
func print_cli() {
	cli := GetCLI()
	clog.Info(fmt.Sprintf("%#v", cli))
}

func Test_get_cli(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	_ = LoadCmd("")
	print_cli()
}

// test load excommand [passed]
func Test_load_excommand(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	_ = LoadCmd("")
	_ = LoadCmd("../../demo/excommand.json")
}

// test get command [passed]
func Test_get_command(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	_ = LoadCmd("")
	verb := "server:start"
	delimiter := ":"
	startCmd := GetCommand(verb, delimiter)
	clog.Info(fmt.Sprintf("%p", startCmd))
}

// test mount flags [passed]
func Test_mount_flags(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	_ = LoadCmd("")
}

// test get flag value [passed]
func Test_get_flagvalue(t *testing.T) {
	config.LoadConfig("../../demo/dev.yaml")
	_ = LoadCmd("")
	// try to get flag velue
	verb := ":server:start"
	delimiter := _default_delimiter
	cmd := GetCommand(verb, delimiter)
	lFlags := cmd.LocalFlags()
	pFlags := cmd.PersistentFlags()
	clog.Info(fmt.Sprintf("%#v", lFlags))
	clog.Info(fmt.Sprintf("%#v", pFlags))
	// show flag value
	clog.Info(fmt.Sprintf("%v", pFlags.Lookup("host")))
	clog.Info(fmt.Sprintf("%v", pFlags.Lookup("port")))
	clog.Info(fmt.Sprintf("%v", lFlags.Lookup("config")))
}
