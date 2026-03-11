package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

func init() {
}

const (
	// define valid flag type in command file.
	FLAG_TYPE_INT    = "int"
	FLAG_TYPE_FLOAT  = "float"
	FLAG_TYPE_STRING = "string"
	FLAG_TYPE_BOOL   = "bool"
)

var (
	_default_command_path = "../../command.json"
	// _dict_command         config.DataDict[any]
	_default_delimiter = ":"
)

type (
	// Flag store flag information.
	Flag struct {
		FullName  string `json:"fullName"`  //  full name
		ShortName string `json:"shortName"` //  short name
		Type      string `json:"type"`      //  value type
		Desc      string `json:"desc"`      // description
		Default   string `json:"default"`   // default value
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

// DefaultCommandPath set default command path
func DefaultCommandPath(path string) {
	_default_command_path = path
}

func DefaultDelimiter(delimiter string) {
	_default_delimiter = delimiter
}

// LoadCmd parse the specified file to the cli instance and
// build the instruction dictionary.
func LoadCmd(path string) *Cli {
	if path == "" {
		path = _default_command_path
	}
	var cli Cli
	// todo: Perhaps different configuration file suffixes could be used to parse this,
	// but JSON is already sufficient for the current functionality.
	if err := util.ParseJsonFile(path, &cli); err != nil {
		clog.Panic(err.Error())
		return nil
	}
	// put cli into config dict
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	configDict.Record(config.DATAKEY_CLI, &cli)
	// load all command to dict_key(_dict_command) data dict
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
func GetCLI() *Cli {
	configDict := config.GetDict(config.DICTKEY_CONFIG)
	cli, ok := configDict.Find(config.DATAKEY_CLI).Value().(*Cli)
	if !ok {
		clog.Panic(fmt.Sprintf("from data(%s) can't assert to type(*CLi)", palette.Red(config.DATAKEY_CLI)))
	}
	return cli
}

// GetCommand return the cobra's Command from command dict and will panic if not exists the specific verb.
// It will set internal command dict cache in package cli.
func GetCommand(verb string, delimiter string) *cobra.Command {
	// todo: delimiter == ["-", "_"]?
	cmdKey := verb
	commandDict := config.GetDict(config.DICTKEY_COMMAND)
	ccmd, ok := commandDict.Find(cmdKey).Value().(*cobra.Command)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*cobra.Command) fail", palette.Red(cmdKey)))
	}
	return ccmd
}

// MountCmd performs the same operation as the built-in mountCmd,
// except that it mounts the commands to the specified dictionary.
// If the specified dictionary does not exist, a panic is triggered.
func MountCmd(verb string, cmd *Command, dictkey string) *cobra.Command {
	dict := config.GetDict(config.DictKey(dictkey))
	return mountCmd(verb, cmd, dict)
}

// The idea here is to try to enter from any top-level instruction in the already
// successfully parsed CLI structure, but obviously this can be automatically
// obtained during the recursion process, so it is retained here, but there is a
// fast flag loading version.
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
	// mountFlag(verb, curCmd) -- Old version implementation
	quickMountFLag(verb, cmd, curCmd)
	dict.Record(verb, curCmd)
	for i := range cmd.SubCommand {
		nextCmd := mountCmd(verb, &cmd.SubCommand[i], dict)
		if nextCmd != nil {
			curCmd.AddCommand(nextCmd)
		}
	}
	return curCmd
}

// using recursive search
func mountFlag(verb string, cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cli := GetCLI()
	flags := cmd.Flags()
	for _, f := range ParseLocalFlags(verb, _default_delimiter, cli) {
		expandFlag(flags, f)
	}
	for _, f := range ParsePersistenFlags(verb, _default_delimiter, cli) {
		expandFlag(flags, f)
	}
	clog.Info(fmt.Sprintf("for verb(%s) mount local and persistent flags", palette.SkyBlue(verb)))
}

// using current context information
func quickMountFLag(verb string, originCmd *Command, treeCmd *cobra.Command) {
	if originCmd == nil || treeCmd == nil {
		return
	}
	flags := treeCmd.Flags()
	for _, f := range originCmd.LocalFlags {
		expandFlag(flags, f)
	}
	clog.Info(fmt.Sprintf("parse verb(%s) local +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(originCmd.LocalFlags))))
	for _, f := range originCmd.PersistentFlags {
		expandFlag(flags, f)
	}
	clog.Info(fmt.Sprintf("parse verb(%s) persistent +[%s] Flags", palette.SkyBlue(verb), palette.Green(len(originCmd.LocalFlags))))
	clog.Info(fmt.Sprintf("for verb(%s) mount local and persistent flags", palette.SkyBlue(verb)))
}

func expandFlag(fset *pflag.FlagSet, f Flag) {
	if fset.Lookup(f.FullName) != nil {
		return
	}
	if fset.ShorthandLookup(f.ShortName) != nil {
		clog.Panic(fmt.Sprintf("Shorthand '%s' is already used, skipping shorthand for %s", f.ShortName, f.FullName))
		f.ShortName = ""
	}
	switch f.Type {
	case FLAG_TYPE_BOOL:
		v, err := strconv.ParseBool(f.Default)
		if err != nil {
			clog.Error(fmt.Sprintf("flag value(%s) mismatch type(%s)", palette.Red(f.Default), palette.Red(FLAG_TYPE_BOOL)))
		}
		fset.BoolP(f.FullName, f.ShortName, v, f.Desc)
	case FLAG_TYPE_INT:
		v, err := strconv.Atoi(f.Default)
		if err != nil {
			clog.Error(fmt.Sprintf("flag value(%s) mismatch type(%s)", palette.Red(f.Default), palette.Red(FLAG_TYPE_INT)))
		}
		fset.IntP(f.FullName, f.ShortName, v, f.Desc)
	case FLAG_TYPE_FLOAT:
		v, err := strconv.ParseFloat(f.Default, 1<<6)
		if err != nil {
			clog.Error(fmt.Sprintf("flag value(%s) mismatch type(%s)", palette.Red(f.Default), palette.Red(FLAG_TYPE_FLOAT)))
		}
		fset.Float64P(f.FullName, f.ShortName, v, f.Desc)
	case FLAG_TYPE_STRING:
		fset.StringP(f.FullName, f.ShortName, f.Default, f.Desc)
	default:
		clog.Warn(fmt.Sprintf("for flag(%s) with invalid type(%s)", palette.SkyBlue(f.FullName), palette.Red(f.Type)))
	}
}
