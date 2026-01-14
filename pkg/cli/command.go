package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

const (
	_default_command_path = "../../command.json"

	_flag_type_int    = "int"
	_flag_type_float  = "float"
	_flag_type_string = "string"
	_flag_type_bool   = "bool"

	_default_value_int    = 0
	_default_value_float  = 0
	_default_value_string = ""
	_default_value_bool   = false
)

var (
	_dict_command *config.DataDict[any]
)

type (
	// Flag store flag information.
	Flag struct {
		FullName  string `json:"fullName"`  //  full name
		ShortName string `json:"shortName"` //  short name
		Type      string `json:"type"`      //  value type
		Desc      string `json:"desc"`      // description
	}
	// Command store command information.
	Command struct {
		Verb            string    `json:"verb"`            // called verb
		ShortDesc       string    `json:"shortDesc"`       // short description for verb
		LongDesc        string    `json:"longDesc"`        // long description for verb
		PersistentFlags []Flag    `json:"persistentFlags"` // persistent flags for current command
		LocalFlags      []Flag    `json:"localFlags"`      // local flags for current command
		SubCommand      []Command `json:"subCommands"`     // subcommands for current command
	}
	// Cli is the standard format after parsing the command file.
	Cli struct {
		App      string    `json:"app"`      // cli name used for root command
		Entry    []string  `json:"entry"`    // execute entry, like Package main.
		Version  string    `json:"version"`  // cli version
		Intro    string    `json:"intro"`    // cli introduction
		Commands []Command `json:"commands"` // all cli commands
	}
)

// LoadCmd parse the specified file to the cli instance and
// build the instruction dictionary.
func LoadCmd(path string) *Cli {
	if path == "" {
		path = _default_command_path
	}
	var cli Cli
	if err := util.ParseJsonFile(path, &cli); err != nil {
		clog.Panic(err.Error())
		return nil
	}
	// set default version for cli
	_default_intro = cli.Intro
	_default_version = cli.Version
	// put cli into config dict
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	configDict.Record(config.DATAKEY_CLI, &cli)
	// load all command to dict_key(_dict_command_) data dict
	var cmdDict config.DataDict[any]
	if !config.HasDict(config.DICTKEY_COMMAND) {
		cmdDict = config.NewDataDict[any](config.DICTKEY_COMMAND)
		config.PutDict(cmdDict.Name(), cmdDict)
	} else {
		cmdDict = config.GetDict(config.DICTKEY_COMMAND)
	}
	for i := range cli.Commands {
		// mount command and flags
		_ = mountCmd("", &cli.Commands[i], cmdDict)
	}
	return &cli
}

// GetCmd return the pointer to Cli from config dict and will panic if not exists Cli.
func GetCmd() *Cli {
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	cli, ok := configDict.Find(config.DATAKEY_CLI).Value().(*Cli)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*CLi) fail", palette.Red(config.DATAKEY_CLI)))
	}
	return cli
}

// GetCommand return the cobra's Command from command dict and will panic if not exists the specific verb.
// It will set internal command dict cache in package cli.
func GetCommand(verb string, delimiter string) *cobra.Command {
	// todo: delimiter == ["-", "_"]?
	cmdKey := verb
	if _dict_command == nil {
		_dict_command = new(config.DataDict[any])
		*_dict_command = config.GetDict(config.DICTKEY_COMMAND)
	}
	ccmd, ok := _dict_command.Find(cmdKey).Value().(*cobra.Command)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*cobra.Command) fail", palette.Red(cmdKey)))
	}
	return ccmd
}

// ParsePersistenFlags return the list for persistent flag of specific verb. No panic here.
func ParsePersistenFlags(verb string, delimiter string, entry *Cli) []Flag {
	if entry == nil {
		clog.Error(fmt.Sprintf("from invalid cli entry find verb(%s) persistent flags", palette.Red(verb)))
		return nil
	}
	for i := range entry.Commands {
		if cmd, find := findCommand(verb, delimiter, entry.Commands[i]); find {
			clog.Info(fmt.Sprintf("parse verb(%s) persistent +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(cmd.PersistentFlags))))
			return cmd.PersistentFlags
		}
	}
	clog.Warn(fmt.Sprintf("not exists verb(%s)", palette.Red(verb)))
	return nil
}

// ParseLocalFlags return the list for local flag of specific verb. No panic here.
func ParseLocalFlags(verb string, delimiter string, entry *Cli) []Flag {
	if entry == nil {
		clog.Error(fmt.Sprintf("from invalid cli entry find verb(%s) local flags", palette.Red(verb)))
		return nil
	}
	for i := range entry.Commands {
		if cmd, find := findCommand(verb, delimiter, entry.Commands[i]); find {
			clog.Info(fmt.Sprintf("parse verb(%s) local +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(cmd.LocalFlags))))
			return cmd.LocalFlags
		}
	}
	clog.Warn(fmt.Sprintf("not exists verb(%s)", palette.Red(verb)))
	return nil
}

// Only used when initializing flags.
// verb should be like :server:start
func findCommand(verb string, delimiter string, cmd Command) (Command, bool) {
	// remove the first delimiter
	if strings.Index(verb, delimiter) != 0 {
		return Command{}, false
	}
	verb = strings.Replace(verb, delimiter, "", 1)
	if verb == cmd.Verb {
		return cmd, true
	}
	// verb should shrink prefix if [prefix]:[suffix] and prefix just equal current verb.
	// otherwise add delimiter to prefix.
	idx := strings.Index(verb, delimiter)
	if idx != -1 && string(verb[:idx]) == cmd.Verb {
		// server:start => :start
		verb = verb[idx:]
	} else {
		// server:start => :server:start
		verb = delimiter + verb
	}
	var c Command
	find := false
	for i := range cmd.SubCommand {
		c, find = findCommand(verb, delimiter, cmd.SubCommand[i])
		if find {
			break
		}
	}
	return c, find
}

func mountCmd(verb string, cmd *Command, dict config.DataDict[any]) *cobra.Command {
	if cmd == nil {
		return nil
	}
	verb += _default_delimiter + cmd.Verb
	curCmd := &cobra.Command{
		Use:   cmd.Verb,
		Short: cmd.ShortDesc,
		Long:  cmd.LongDesc,
	}
	mountFlag(verb, curCmd)
	dict.Record(verb, curCmd)
	for i := range cmd.SubCommand {
		nextCmd := mountCmd(verb, &cmd.SubCommand[i], dict)
		if nextCmd != nil {
			curCmd.AddCommand(nextCmd)
		}
	}
	return curCmd
}

func mountFlag(verb string, cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cli := GetCmd()
	lFlags := cmd.LocalFlags()
	pFlags := cmd.PersistentFlags()
	for _, f := range ParseLocalFlags(verb, _default_delimiter, cli) {
		expandFlag(lFlags, f)
	}
	for _, f := range ParsePersistenFlags(verb, _default_delimiter, cli) {
		expandFlag(pFlags, f)
	}
	clog.Info(fmt.Sprintf("for verb(%s) mount local and persistent flags", palette.SkyBlue(verb)))
}

func expandFlag(fset *pflag.FlagSet, f Flag) {
	switch f.Type {
	case _flag_type_bool:
		boolVar := new(bool)
		fset.BoolVarP(boolVar, f.FullName, f.ShortName, _default_value_bool, f.Desc)
	case _flag_type_int:
		intVar := new(int)
		fset.IntVarP(intVar, f.FullName, f.ShortName, _default_value_int, f.Desc)
	case _flag_type_float:
		floatVar := new(float64)
		fset.Float64VarP(floatVar, f.FullName, f.ShortName, _default_value_float, f.Desc)
	case _flag_type_string:
		stringVar := new(string)
		fset.StringVarP(stringVar, f.FullName, f.ShortName, _default_value_string, f.Desc)
	default:
		clog.Warn(fmt.Sprintf("for flag(%s) with invalid type(%s)", palette.SkyBlue(f.FullName), palette.Red(f.Type)))
	}
}
