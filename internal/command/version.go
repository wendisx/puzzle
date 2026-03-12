package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/cli"
	"github.com/wendisx/puzzle/pkg/config"
)

const (
	_artfont = ` 
┌─┐┬ ┬┌─┐┌─┐┬  ┌─┐
├─┘│ │┌─┘┌─┘│  ├┤ 
┴  └─┘└─┘└─┘┴─┘└─┘
  intro: %s
version: %s
  build: %s
   desc: %s
   date: %s
`
)

var (
	_verb_version  = "version"
	_short_version = "Show puzzle version info"
	_long_version  = ""

	INTRO   = "enjoy using puzzle! :-)"
	VERSION = "unknown"
	BUILD   = "unknown"
	DESC    = "a cli tool to simplifies the generation of project."
	DATE    = "-"
)

// MountVersion mount the verb `-version` to the command tree.
// For details, see the `command.json` file in the project structure to
// find the structure of the corresponding command.
func MountBuiltinVersion(rootCmd *cobra.Command) {
	_versionCmd := &cli.Command{
		Verb:      _verb_version,
		ShortDesc: _short_version,
		LongDesc:  _long_version,
	}
	versionCmd := cli.MountCmd("", _versionCmd, config.DICTKEY_COMMAND)
	versionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(os.Stderr, _artfont, INTRO, VERSION, BUILD, DESC, DATE)
		return nil
	}
	rootCmd.AddCommand(versionCmd)
}
