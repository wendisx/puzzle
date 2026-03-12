package command

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/wendisx/puzzle/pkg/config"
)

const (
	_temptype_file uint8 = iota + 1
	_temptype_dir

	_api_host = `https://api.github.com`
	_raw_host = `https://raw.githubusercontent.com`

	_gitapi_content_info = `/repos/%s/%s/contents/%s?ref=%s` // /owner/repo/path/ref[branch|tag|hash]
	_gitapi_raw_data     = `/%s/%s/%s/%s`                    // /owner/repo/branch/path

	_type_file = "file"
	_type_dir  = "dir"
)

var (
	_template_map = map[string]templateSrc{
		"_test_new":  templateSrc{_temptype_file, "/readme.md", "md"},
		"command":    templateSrc{_temptype_file, "/template/new/command.json", "json"},
		"config":     templateSrc{_temptype_file, "/template/new/config.yaml", "yaml"},
		"config_dev": templateSrc{_temptype_file, "/template/new/config-dev.yaml", "yaml"},
		"_test_init": templateSrc{_temptype_dir, "/pkg", ""},
		"echo":       templateSrc{_temptype_dir, "/template/init/puzzle-echo", ""},
		"chi":        templateSrc{_temptype_dir, "/template/init/puzzle-chi", ""},
		"gin":        templateSrc{_temptype_dir, "/template/init/puzzle-gin", ""},
	}
	// file pull pool
	FileBufPool = sync.Pool{
		New: func() any {
			return new(bufio.Writer)
		},
	}
	// file pull synchronizer
	FileSynchronizer sync.WaitGroup
	FileErrChan      = make(chan error)
)

type (
	// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28&versionId=free-pro-team%40latest&restPage=using-the-rest-api-to-interact-with-your-git-database&category=repos&subcategory=webhooks
	ContentNode struct {
		Name    string        `json:"name"`
		Path    string        `json:"path"`
		Size    int           `json:"size"`
		Type    string        `json:"type"`
		Url     string        `json:"url"`
		Entries []ContentNode `json:"entries"`
	}
	GithubApi struct{}
)

func (a GithubApi) RawFile(format string, remote *config.GithubConfig, path string) string {
	uri := fmt.Sprintf(format, remote.UserName, remote.Repos[remote.ActiveRepo].Name, remote.Repos[remote.ActiveRepo].Ref, path)
	if remote.RawHost == "" {
		remote.RawHost = _raw_host
	}
	return remote.RawHost + filepath.Clean(uri)
}

func (a GithubApi) PathInfo(format string, remote *config.GithubConfig, path string) string {
	uri := fmt.Sprintf(format, remote.UserName, remote.Repos[remote.ActiveRepo].Name, path, remote.Repos[remote.ActiveRepo].Ref)
	if remote.ApiHost == "" {
		remote.ApiHost = _raw_host
	}
	return remote.ApiHost + filepath.Clean(uri)
}

func GetFileBuf(w io.Writer) *bufio.Writer {
	buf := FileBufPool.Get().(*bufio.Writer)
	buf.Reset(w)
	return buf
}

func PutFileBuf(buf *bufio.Writer) {
	buf.Reset(nil)
	FileBufPool.Put(buf)
}
