package app

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

const (
	// PerPage ...
	PerPage = 99999
)

// GitLabClient ...
type GitLabClient interface {
	Namespaces() []*gitlab.Namespace
	Projects() []*gitlab.Project
}

type gitLabClient struct {
	cli *gitlab.Client
}

// NewGitLabClient ...
func NewGitLabClient(cfg *GitlabConfig) GitLabClient {
	cli := gitlab.NewClient(nil, cfg.PrivateToken)
	cli.SetBaseURL(fmt.Sprintf("%s/api/v3", cfg.HostURL))
	return &gitLabClient{
		cli: cli,
	}
}

// Namespaces ...
func (c *gitLabClient) Namespaces() []*gitlab.Namespace {
	namespaces, res, err := c.cli.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: PerPage,
		},
	})
	if err != nil {
		fmt.Printf("[Namespaces] %s\n", err.Error()) // TODO
		return nil
	}
	if res.Status != "200 OK" {
		fmt.Printf("[Namespaces] HTTP Response Status:%s\n", res.Status) // TODO
		return nil
	}
	return namespaces
}

// Projects ...
func (c *gitLabClient) Projects() []*gitlab.Project {
	projects, res, err := c.cli.Projects.ListProjects(&gitlab.ListProjectsOptions{
		OrderBy: gitlab.String("name"),
		Sort:    gitlab.String("asc"),
		ListOptions: gitlab.ListOptions{
			PerPage: PerPage,
		},
	})
	if err != nil {
		fmt.Printf("[Projects] %s\n", err.Error()) // TODO
		return nil
	}
	if res.Status != "200 OK" {
		fmt.Printf("[Projects] HTTP Response Status:%s\n", res.Status) // TODO
		return nil
	}
	return projects
}
