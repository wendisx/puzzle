package cli

import (
	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

var (
	_exists_root       = false
	_verb_root         = "root"
	_default_delimiter = ":"

	_default_root_use   = "puzzle"
	_default_root_short = ""
	_default_root_long  = ""
)

func RootVerb(verb string) {
	_verb_root = verb
}

func DefaultDelimiter(delimiter string) {
	_default_delimiter = delimiter
}

func mountRoot() *cobra.Command {
	var rootCmd *cobra.Command
	if !_exists_root {
		rootCmd = &cobra.Command{
			Use:   _default_root_use,
			Short: _default_root_short,
			Long:  _default_root_long,
		}
		if _dict_command == nil {
			_dict_command = new(config.DataDict[any])
			*_dict_command = config.GetDict(config.DICTKEY_COMMAND)
		}
		_dict_command.Record(_verb_root, rootCmd)
		_exists_root = true
	} else {
		rootCmd = GetCommand(_verb_root, _default_delimiter)
	}
	return rootCmd
}

// Instruction execution entry.
// mountFunc represents that the root command is passed into this method,
// which is used to expand the mounting custom command set.
// There are some default internal mount method, like [mountVersion] etc.
func Execute(mountFuncs ...func(*cobra.Command)) {
	rootCmd := mountRoot()
	mountVersion(rootCmd)
	mountServer(rootCmd)
	for i := range mountFuncs {
		mountFuncs[i](rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		clog.Fatal(err.Error())
	}
}
