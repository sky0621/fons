package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	OutputDir string
	Gitlab    *GitlabConfig
	Filter    *FilterConfig
}

// GitlabConfig ...
type GitlabConfig struct {
	HostURL      string
	PrivateToken string
	GitCloneURL  string
	Branch       string
}

// FilterConfig ...
type FilterConfig struct {
	TargetNameSpaces string
	ExcludeProjects  string
}

// NewConfig ...
func NewConfig(path string) (*Config, error) {
	var conf Config

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Host4GitCommand ...
func (c *Config) Host4GitCommand(pathWithNamespace string) string {
	if c == nil {
		return ""
	}
	if c.Gitlab == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s.git", c.Gitlab.GitCloneURL, pathWithNamespace)
}

// ExcludeProjectSlice ...
func (c *Config) ExcludeProjectSlice() []string {
	if c == nil {
		return nil
	}
	if c.Filter == nil {
		return nil
	}
	return strings.Split(c.Filter.ExcludeProjects, ",")
}

// TargetNameSpacesSlice ...
func (c *Config) TargetNameSpacesSlice() []string {
	if c == nil {
		return nil
	}
	if c.Filter == nil {
		return nil
	}
	return strings.Split(c.Filter.TargetNameSpaces, ",")
}

// IsTargetNamespace ...
func (c *Config) IsTargetNamespace(path string) bool {
	targetNameSpacesSlice := c.TargetNameSpacesSlice()
	if targetNameSpacesSlice == nil {
		return false
	}
	for _, targetNameSpace := range targetNameSpacesSlice {
		if strings.Contains(path, targetNameSpace) {
			return true
		}
	}
	return false
}

// TargetProjectPathInfos ...
func (c *Config) TargetProjectPathInfos(path string) []os.FileInfo {
	if c == nil {
		return nil
	}
	_, err := os.Stat(filepath.Join(c.OutputDir, path))
	if err != nil {
		err = os.Mkdir(filepath.Join(c.OutputDir, path), 0777)
		if err != nil {
			panic(err) // TODO
		}
	}

	files, err := ioutil.ReadDir(filepath.Join(c.OutputDir, path))
	if err != nil {
		panic(err) // TODO
	}

	return files
}

// IsExcludeProject ...
func (c *Config) IsExcludeProject(p string) bool {
	for _, outPrj := range c.ExcludeProjectSlice() {
		if strings.Contains(p, outPrj) {
			return true
		}
	}
	return false
}
