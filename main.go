package main

import (
	"errors"
	"fmt"

	"io/ioutil"

	"os"

	"os/exec"

	"path/filepath"

	"flag"

	"strings"

	"github.com/xanzy/go-gitlab"
)

const (
	perPage = 99999
)

// TODO 機能実現スピード最優先での実装なので要リファクタ
func main() {
	fmt.Println("Start")

	f := flag.String("f", "./config.toml", "Config File")
	flag.Parse()

	ReadConfig(*f)
	cfg := NewConfig()

	gitCli := gitlab.NewClient(nil, cfg.PrivateToken)
	gitCli.SetBaseURL(fmt.Sprintf("%s/api/v3", cfg.GitlabApiUrl))

	namespaces, res, err := gitCli.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{
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

	filterOutProjects := strings.Split(cfg.FilterOutProject, ",")

	for _, ns := range namespaces {
		if !strings.Contains(ns.Path, cfg.FilterInNameSpace) {
			continue
		}
		projects, res, err := gitCli.Projects.ListProjects(&gitlab.ListProjectsOptions{
			OrderBy: gitlab.String("name"),
			Sort:    gitlab.String("asc"),
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

		if len(projects) == 0 {
			continue
		}

		_, err = os.Stat(filepath.Join(cfg.Outdir, ns.Path))
		if err != nil {
			err = os.Mkdir(filepath.Join(cfg.Outdir, ns.Path), 0777)
			if err != nil {
				panic(err)
			}
		}

		files, err := ioutil.ReadDir(filepath.Join(cfg.Outdir, ns.Path))
		if err != nil {
			panic(err)
		}

		for _, p := range projects {
			if ns.Path != p.Namespace.Path {
				continue
			}
			if isOutProject(p.Path, filterOutProjects) {
				continue
			}
			fmt.Println(p.PathWithNamespace)
			if exists(files, func(filename string) bool {
				return filename == p.Path
			}) {
				err := os.Chdir(filepath.Join(cfg.Outdir, ns.Path, p.Path))
				if err != nil {
					panic(err)
				}

				cmd := exec.Command("git", "pull")
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			} else {
				cmd := exec.Command("git", "clone", cfg.Host4GitCommand(p.PathWithNamespace), filepath.Join(cfg.Outdir, ns.Path, p.Path))
				err := cmd.Run()
				if err != nil {
					panic(err)
				}

				err = os.Chdir(filepath.Join(cfg.Outdir, ns.Path, p.Path))
				if err != nil {
					panic(err)
				}

				cmd3 := exec.Command("git", "checkout", "-b", cfg.Branch, "origin/"+cfg.Branch)
				err = cmd3.Run()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func isOutProject(p string, outProjects []string) bool {
	for _, outPrj := range outProjects {
		if strings.Contains(p, outPrj) {
			return true
		}
	}
	return false
}

func exists(files []os.FileInfo, fn func(filename string) bool) bool {
	for _, file := range files {
		if exists := fn(file.Name()); exists {
			return true
		}
	}
	return false
}
