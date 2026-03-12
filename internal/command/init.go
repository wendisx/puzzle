package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/wendisx/puzzle/pkg/cli"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
)

/*
	puzzle init -- 初始化工程模板
	init和new的差异: init是脚手架级别的, new只负责单个文件模板的初始化.
	init应该做到:
	1. 指定脚手架的模板从github拉取指定模板到本地指定dest位置, 自动携带依赖图, 但是需要自行本地安装.
	2. init只做初始化工作, 并始终保持最小化构建.
	3. init不尝试检查和做出任何关键性决定, 这些依赖开发者.
	4. 并发拉取工程模板文件.
*/

const (
	_verb_init  = "init"
	_short_init = "initialize project template"
	_long_init  = ""
)

func MountBuiltinInit(rootCmd *cobra.Command) {
	_initCmd := &cli.Command{
		Verb:      _verb_init,
		ShortDesc: _short_init,
		LongDesc:  _long_init,
		LocalFlags: []cli.Flag{
			cli.Flag{"template", "t", cli.FLAG_TYPE_STRING, "Specify project template, like echo, chi, gin", ""},
			cli.Flag{"dest", "d", cli.FLAG_TYPE_STRING, "Specify the location where the template is created.", "."},
		},
	}
	initCmd := cli.MountCmd("", _initCmd, config.DICTKEY_COMMAND)
	initCmd.RunE = func(cmd *cobra.Command, args []string) error {
		tempf, err := cmd.Flags().GetString(_initCmd.LocalFlags[0].FullName)
		dest, err := cmd.Flags().GetString(_initCmd.LocalFlags[1].FullName)
		if err != nil {
			return err
		}
		if _template_map[tempf].Ref == "" {
			return fmt.Errorf("template should be non-empty and valid.")
		}
		f, err := os.Stat(dest)
		if err != nil {
			return err
		}
		if !f.IsDir() {
			dest = "."
		}
		dest = filepath.Clean(dest)
		// align initial directory
		// apiUrl like: https://api.github.com/repos/owner/repo_name/contents/src?ref=branch_name
		remote := config.GetConfig().GithubConfig
		if remote.ApiHost == "" {
			remote.ApiHost = `https://api.github.com`
		}
		src := fmt.Sprintf(_gitapi_content_info,
			remote.UserName,
			remote.Repos[remote.ActiveRepo].Name,
			_template_map[tempf].Ref,
			remote.Repos[remote.ActiveRepo].Ref,
		)
		src = remote.ApiHost + filepath.Clean(src)
		clog.Warn(fmt.Sprint(src))
		// generate a http client
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		clog.Warn("enter syncTemplate")
		err = syncTemplate(&remote, client, src, dest)
		if err != nil {
			_ = os.RemoveAll(dest)
			return err
		}
		select {
		case err = <-FileErrChan:
			return err
		default:
			FileSynchronizer.Wait()
		}
		fmt.Fprintf(os.Stderr, "Sync all template files to %s.\n", dest)
		return nil
	}
	rootCmd.AddCommand(initCmd)
}

func syncTemplate(remote *config.GithubConfig, client *http.Client, src, dest string) error {
	req, err := http.NewRequest(http.MethodGet, src, nil)
	if err != nil {
		return err
	}
	// set header
	req.Header.Set("Accept", "application/vnd.github.object")
	req.Header.Set("Authorization", remote.AccessToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api error: %d", resp.StatusCode)
	}
	// marshal response body to json
	var body ContentNode
	if err = json.NewDecoder(bufio.NewReader(resp.Body)).Decode(&body); err != nil {
		return err
	}
	resp.Body.Close()
	clog.Debug(fmt.Sprintf("%+v", body))
	switch body.Type {
	case _type_dir:
		for i := range body.Entries {
			if body.Entries[i].Type == _type_dir {
				subDest := dest + "/" + body.Entries[i].Name
				err = os.MkdirAll(filepath.Clean(subDest), 0755)
				if err != nil {
					return err
				}
				if err = syncTemplate(remote, client, body.Entries[i].Url, subDest); err != nil {
					return err
				}
			} else {
				FileSynchronizer.Add(1)
				if err = syncTemplate(remote, client, body.Entries[i].Url, dest); err != nil {
					return err
				}
			}
		}
	case _type_file:
		filePath := filepath.Clean(dest + "/" + body.Name)
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}
		go func(filePath string) {
			defer FileSynchronizer.Done()
			f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				FileErrChan <- err
			}
			api := GithubApi{}
			resp, err := client.Get(api.RawFile(_gitapi_raw_data, remote, body.Path))
			if err != nil {
				FileErrChan <- err
			}
			defer resp.Body.Close()
			buf := GetFileBuf(f)
			if _, err = io.Copy(buf, resp.Body); err != nil {
				FileErrChan <- err
			}
			buf.Flush()
			f.Close()
		}(filePath)
	}
	return nil
}
