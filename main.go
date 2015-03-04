package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	Path  string
	Rules map[string]string
}

func main() {
	config := &Config{}
	configData, err := ioutil.ReadFile("rename_rule.yml")
	if err != nil {
		fmt.Println("configuration file read error. must be exists rename_rule.yml.")
		os.Exit(1)
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		fmt.Println("configuration file is not valid format.")
		os.Exit(1)
	}

	file, err := os.Stat(config.Path)
	err = config.RenameFileRecursive(file, config.Path)
	if err != nil {
		fmt.Println("rename error: %s", err)
		os.Exit(1)
	}
}

func (c *Config) RenameFileRecursive(f os.FileInfo, path string) error {
	if f.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			joinedPath := filepath.Join(path, file.Name())
			if err != nil {
				return err
			}
			err := c.RenameFileRecursive(file, joinedPath)
			if err != nil {
				return err
			}
		}
	}

	for pattern, replacement := range c.Rules {
		if ok, err := filepath.Match(pattern, f.Name()); ok && err == nil {
			asteriskToRegexp := regexp.MustCompile("(.*)\\*")
			patternRegexpStr := asteriskToRegexp.ReplaceAllString(pattern, "$1(.*)")
			patternRegexp := regexp.MustCompile(patternRegexpStr)
			holderedReplacement := asteriskToRegexp.ReplaceAllString(replacement, "$1")
			replacedPath := patternRegexp.ReplaceAllString(path, holderedReplacement+"$1")
			fmt.Printf("renamed: %s => %s\n", path, replacedPath)
			err := os.Rename(path, replacedPath)
			if err != nil {
				return err
			}
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
