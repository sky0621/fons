package app

import (
	"fmt"
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
