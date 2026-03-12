// Package config load all configurations and environment variable to program.
package config

import (
	"fmt"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/util"
)

func init() {
	// The dictionary directory is initialized synchronously during
	// configuration package initialization.
	if _dict_directory == nil {
		_dict_directory = new(DictDirectory)
	}
}

const (
	/* dictionary key to find dictionary */
	DICTKEY_CONFIG  = "_dict_config"
	DICTKEY_COMMAND = "_dict_command"
	DICTKEY_SERVER  = "_dict_server"
	DICTKEY_SERVICE = "_dict_service"
	DICTKEY_CLIENT  = "_dict_client"
)

const (
	/* data key to fing data */
	DATAKEY_CONFIG           = "_data_config"
	DATAKEY_ENV              = "_data_env"
	DATAKEY_CLI              = "_data_cli"
	DATAKEY_SERVER_ECHO      = "_data_server_echo"
	DATAKEY_PERMISSION_TABLE = "_data_permission_table"
	DATAKEY_PERMISSION_USER  = "_data_permission_user"
	DATAKEY_DB_REDIS         = "_data_db_redis"
)

var (
	// default config file path
	_default_config_path = "./vendor/puzzle/demo/example.yaml"
)

// Config record record all possible configuration items.
type Config struct {
	EnvConfig    []string     `yaml:"environment" json:"environment"` // special environment config
	GithubConfig GithubConfig `yaml:"github" json:"github"`           // github config
	DBConfig     DBConfig     `yaml:"database" json:"database"`       // database config
	ServerConfig ServerConfig `yaml:"server" json:"server"`           // server config
	SwagConfig   SwagConfig   `yaml:"swagger" json:"swagger"`         // swagger config
}

// DefaultConfigFile set default config file.
func DefaultConfigPath(path string) {
	_default_config_path = path
}

// LoadConfig return a pointer to all config and will panic if not exists the file path.
func LoadConfig(path string) *Config {
	// init config dict here.
	var configDict DataDict[any]
	if HasDict(DICTKEY_CONFIG) {
		configDict = GetDict(DICTKEY_CONFIG)
	} else {
		configDict = NewDataDict[any](DICTKEY_CONFIG)
		PutDict(configDict.Name(), configDict)
	}
	c := &Config{
		DBConfig:     initDBConfig(),
		ServerConfig: initServerConfig(),
		GithubConfig: initGithubConfig(),
	}
	if path == "" {
		path = _default_config_path
	}
	if err := util.ParseYamlFile(path, c); err != nil {
		clog.Panic(fmt.Sprintf("%s", err.Error()))
		return c
	}
	// put Config into data dict
	configDict.Record(DATAKEY_CONFIG, c)
	return c
}

// GetConfig try to get pointer to global config and will panic if not exists the config.
func GetConfig() *Config {
	configDict := GetDict(DICTKEY_CONFIG)
	c, ok := configDict.Find(DATAKEY_CONFIG).Value().(*Config)
	if !ok {
		clog.Panic(fmt.Sprintf("from data(%s) can't assert to type(*Config)", palette.Red(DATAKEY_CONFIG)))
	}
	return c
}
