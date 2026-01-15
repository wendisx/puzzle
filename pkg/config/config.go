// Package config load all configurations and environment variable to program.
package config

import (
	"fmt"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

const (
	// some default key from internal, exactly I'd like to got them from some .json files. :)
	_dict_capacity = 1 << 10
)

const (
	/* dictionary key to find dictionary */
	DICTKEY_CONFIG  = "_dict_config_"
	DICTKEY_COMMAND = "_dict_command"
)

const (
	/* data key to fing data */
	DATAKEY_CONFIG = "_data_config_"
	DATAKEY_ENV    = "_data_env_"
	DATAKEY_CLI    = "_data_cli_"
)

var (
	_default_config_file = "./vendor/puzzle/demo/example.yaml"
)

// Config record record all possible configuration items.
type Config struct {
	EnvConfig    []string     `yaml:"environment"` // special environment config
	DBConfig     DBConfig     `yaml:"database"`    // database config
	ServerConfig ServerConfig `yaml:"server"`      // server config
	SwagConfig   SwagConfig   `yaml:"swagger"`     // swagger config
}

func init() {
	// init dict directory
	_dict_directory = new(DictDirectory)
	// init config dict here.
	configDict := NewDataDict[any](DICTKEY_CONFIG)
	PutDict(configDict.Name(), configDict)
}

// DefaultConfigFile set global default config file.
func DefaultConfigFile(path string) {
	_default_config_file = path
}

// LoadConfig return a pointer to all config and will panic if not exists the file path.
func LoadConfig(path string) *Config {
	c := &Config{
		DBConfig:     initDBConfig(),
		ServerConfig: initServerConfig(),
	}
	if path == "" {
		path = _default_config_file
	}
	if err := util.ParseYamlFile(path, c); err != nil {
		clog.Panic(fmt.Sprintf("%s", err.Error()))
		return c
	}
	// put Config into data dict
	configDict := GetDict(DICTKEY_CONFIG)
	configDict.Record(DATAKEY_CONFIG, c)
	return c
}

// GetConfig try to get pointer to global config and will panic if not exists the config.
func GetConfig() *Config {
	configDict := GetDict(DICTKEY_CONFIG)
	c, ok := configDict.Find(DATAKEY_CONFIG).Value().(*Config)
	if !ok {
		clog.Panic(fmt.Sprintf("from data_key(%s) assert to type(*Config) fail", palette.Red(DATAKEY_CONFIG)))
	}
	return c
}
