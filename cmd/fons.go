package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sky0621/fons/app"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	perPage = 99999
)

var (
	configPath = flag.String("f", "./config.toml", "Config File")
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()
	cfg, err := app.NewConfig(*configPath)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	glCli := app.NewGitLabClient(cfg.Gitlab)
	namespaces, res, err := glCli.Namespaces(&gitlab.ListNamespacesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: perPage,
		},
	})
	if err != nil {
		panic(err)
	}
	if res.Status != "200 OK" {
		panic(errors.New("not 200 OK"))
	}

	// FIXME goroutine
	fmt.Println(namespaces)

	return 0
}
