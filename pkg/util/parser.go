package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"go.yaml.in/yaml/v3"
)

func ParseJsonString(sr string, dest any) error {
	if err := json.Unmarshal([]byte(sr), dest); err != nil {
		return err
	}
	clog.Info(fmt.Sprintf("parse %s successfully", palette.SkyBlue(sr)))
	return nil
}

func ParseJsonFile(path string, dest any) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(dest); err != nil {
		return err
	}
	clog.Info(fmt.Sprintf("parse %s successfully", palette.SkyBlue(absPath)))
	return nil
}

func ParseYamlString(sr string, dest any) error {
	if err := yaml.Unmarshal([]byte(sr), dest); err != nil {
		return err
	}
	clog.Info(fmt.Sprintf("parse %s successfully", palette.SkyBlue(sr)))
	return nil
}

func ParseYamlFile(path string, dest any) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = yaml.NewDecoder(bufio.NewReader(f)).Decode(dest); err != nil {
		return err
	}
	clog.Info(fmt.Sprintf("parse %s successfully", palette.SkyBlue(absPath)))
	return nil
}
