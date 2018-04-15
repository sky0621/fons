package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	GitlabApiUrl      string
	PrivateToken      string
	GitlabSshUrl      string
	Branch            string
	Outdir            string
	FilterInNameSpace string
	FilterOutProject  string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		GitlabApiUrl:      viper.GetString("gitlabApiUrl"),
		PrivateToken:      viper.GetString("privateToken"),
		GitlabSshUrl:      viper.GetString("gitlabSshUrl"),
		Branch:            viper.GetString("branch"),
		Outdir:            viper.GetString("outdir"),
		FilterInNameSpace: viper.GetString("filterInNameSpace"),
		FilterOutProject:  viper.GetString("filterOutProject"),
	}
}

// ReadConfig ...
func ReadConfig(configFilePath string) error {
	viper.SetConfigFile(configFilePath)
	return viper.ReadInConfig()
}

// Host4GitCommand ...
func (c *Config) Host4GitCommand(pathWithNamespace string) string {
	return fmt.Sprintf("%s/%s.git", c.GitlabSshUrl, pathWithNamespace)
}
