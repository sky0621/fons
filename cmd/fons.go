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

// TODO 機能実現スピード最優先での実装なので要リファクタ
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

	exitCh := make(chan struct{})
	go func(exitCh chan struct{}) {
		// FIXME goroutine
		fmt.Println(namespaces)

	}(exitCh)

	fmt.Println("before exitCh")
	exitCh <- struct{}{}
	fmt.Println("after exitCh")

	return 0
}
