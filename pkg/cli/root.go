package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/config"
)

var (
	_exists_root = false
	_verb_root   = "root"

	_default_root_use   = "puzzle"
	_default_root_short = ""
	_default_root_long  = ""
)

func RootVerb(verb string) {
	_verb_root = verb
}

func RootShort(short string) {
	_default_root_short = short
}

func RootLong(long string) {
	_default_root_long = long
}

func mountRoot() *cobra.Command {
	var rootCmd *cobra.Command
	if !_exists_root {
		rootCmd = &cobra.Command{
			Use:   _default_root_use,
			Short: _default_root_short,
			Long:  _default_root_long,
		}
		commandDict := config.GetDict(config.DICTKEY_COMMAND)
		commandDict.Record(_verb_root, rootCmd)
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
// Perhaps we could obtain the context here to control which commands can be added to the
// command tree, but the simplest way is not to set a default command to add, ensuring that
// the command tree is always empty. This may be useful for some commands. Therefore, we
// can only add commands that are not related to business logic to the command tree as tools,
// and the commands related to writing later can be set in the actual project.
func Execute(mountFuncs ...func(*cobra.Command)) {
	exitText := `Perhaps you need to use the **help** command for some assistance.`
	rootCmd := mountRoot()
	// todo: Change to optional mounting, only mounting general utility commands.
	// MountVersion(rootCmd)
	// MountServer(rootCmd)
	for i := range mountFuncs {
		mountFuncs[i](rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, exitText)
		os.Exit(1)
	}
}
