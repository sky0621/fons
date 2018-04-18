package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sky0621/fons/app"
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

	exitCh := make(chan struct{})
	defer func() {
		close(exitCh)
	}()

	go func(exitCh chan struct{}, cfg *app.Config) {

		maxConcurrentGoroutineNum := runtime.NumCPU() * 3
		fmt.Printf("maxConcurrentGoroutineNum:%d\n", maxConcurrentGoroutineNum)
		semaphore := make(chan struct{}, maxConcurrentGoroutineNum)
		defer func() {
			close(semaphore)
		}()

		glCli := app.NewGitLabClient(cfg.Gitlab)

		// ネームスペース数が膨大になることは想定しないため、同期ループ
		for _, ns := range glCli.Namespaces() {
			fmt.Printf("namespace.Path:%s\n", ns.Path)
			if !cfg.IsTargetNamespace(ns.Path) {
				continue
			}

			pathInfos := cfg.TargetProjectPathInfos(ns.Path)

			for _, project := range glCli.Projects() {
				if ns.Path != project.Namespace.Path {
					continue
				}
				if cfg.IsExcludeProject(project.Path) {
					continue
				}

				semaphore <- struct{}{}

				// TODO 関数化
				go func(semaphore chan struct{}, cfg *app.Config, namespacePath, projectPath string) {
					defer func() {
						<-semaphore
					}()

					if exists(pathInfos, func(filename string) bool {
						return filename == project.Path
					}) {
						err := os.Chdir(filepath.Join(cfg.OutputDir, namespacePath, projectPath))
						if err != nil {
							panic(err)
						}

						cmd := exec.Command("git", "pull")
						err = cmd.Run()
						if err != nil {
							fmt.Println(err)
						}
					} else {
						cmd := exec.Command("git", "clone", cfg.Host4GitCommand(project.PathWithNamespace), filepath.Join(cfg.OutputDir, namespacePath, projectPath))
						err := cmd.Run()
						if err != nil {
							panic(err)
						}

						err = os.Chdir(filepath.Join(cfg.OutputDir, namespacePath, projectPath))
						if err != nil {
							panic(err)
						}

						cmd3 := exec.Command("git", "checkout", "-b", cfg.Gitlab.Branch, "origin/"+cfg.Gitlab.Branch)
						err = cmd3.Run()
						if err != nil {
							fmt.Println(err)
						}
					}
				}(semaphore, cfg, ns.Path, project.Path)
			}
		}
		exitCh <- struct{}{}
	}(exitCh, cfg)

	fmt.Println("before exitCh")
	<-exitCh
	fmt.Println("after exitCh")

	return 0
}

func exists(files []os.FileInfo, fn func(filename string) bool) bool {
	for _, file := range files {
		if exists := fn(file.Name()); exists {
			return true
		}
	}
	return false
}
