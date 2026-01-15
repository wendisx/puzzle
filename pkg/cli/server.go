package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/server"
)

const (
	_flag_handler = "handler"
	_flag_check   = "check"
	_flag_swag    = "swag"
)

var (
	_verb_server       = ":server"
	_verb_server_start = ":server:start"
)

func mountServer(rootCmd *cobra.Command) {
	startCmd := GetCommand(_verb_server_start, _default_delimiter)
	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		hf, err := cmd.Flags().GetString(_flag_handler)
		checkf, err := cmd.Flags().GetBool(_flag_check)
		swagf, err := cmd.Flags().GetBool(_flag_swag)
		if err != nil {
			clog.Error(err.Error())
			return err
		}
		server := server.InitWebServer(hf)
		server.WithCheckRoute(checkf)
		server.WithSwagRoute(swagf)
		server.Start()
		return nil
	}
	serverCmd := GetCommand(_verb_server, _default_delimiter)
	serverCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// todo: show server usage...
		clog.Fatal(fmt.Sprintf("Maybe you need some help with the server directive"))
		return nil
	}
	rootCmd.AddCommand(serverCmd)
}
