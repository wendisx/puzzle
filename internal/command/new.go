package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/cli"
	"github.com/wendisx/puzzle/pkg/config"
)

const (
	_verb_new  = "new"
	_short_new = "generate new template file"
	_long_new  = ""
)

type (
	templateSrc struct {
		Type   uint8 // identify resource types
		Ref    string
		Suffix string
	}
)

func MountBuiltinNew(rootCmd *cobra.Command) {
	_newCmd := &cli.Command{
		Verb:      _verb_new,
		ShortDesc: _short_new,
		LongDesc:  _long_new,
		LocalFlags: []cli.Flag{
			cli.Flag{"template", "t", cli.FLAG_TYPE_STRING, "Specify file template, like command, config, config-dev", ""},
			cli.Flag{"dest", "d", cli.FLAG_TYPE_STRING, "Specify the location where the file is created.", "."},
		},
	}
	newCmd := cli.MountCmd("", _newCmd, config.DICTKEY_COMMAND)
	newCmd.RunE = func(cmd *cobra.Command, args []string) error {
		f_template, err := cmd.Flags().GetString(_newCmd.LocalFlags[0].FullName)
		f_dest, err := cmd.Flags().GetString(_newCmd.LocalFlags[1].FullName)
		if err != nil {
			return err
		}
		return execNew(f_template, f_dest)
	}
	rootCmd.AddCommand(newCmd)
}

func execNew(template, dest string) error {
	if _template_map[template].Ref == "" {
		return errors.New("template should be non-empty and valid.")
	}
	// check dest: dest should be a dir or file with the same suffix.
	info, err := os.Stat(dest)
	if err != nil {
		return err
	}
	fileName := dest + _template_map[template].Ref
	if info.Mode().IsRegular() {
		ext := strings.ToLower(filepath.Ext(dest))
		if ext != _template_map[template].Suffix {
			return fmt.Errorf("mismatched suffix: %s but need %s", ext, _template_map[template].Suffix)
		}
		fileName = dest
	}
	fileName = filepath.Clean(fileName)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	remote := config.GetConfig().GithubConfig
	if remote.RawHost == "" {
		remote.RawHost = `https://raw.githubusercontent.com`
	}
	buf := bufio.NewWriter(f)
	api := GithubApi{}
	resp, err := http.Get(api.RawFile(_gitapi_raw_data, &remote, _template_map[template].Ref))
	if err != nil || resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
		_ = os.Remove(fileName)
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return err
	}
	buf.Flush()
	fmt.Fprintf(os.Stderr, "created: %s", fileName)
	return nil
}
