package app

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GitLabClient ...
type GitLabClient interface {
	Namespaces(opt *gitlab.ListNamespacesOptions, options ...gitlab.OptionFunc) ([]*gitlab.Namespace, *gitlab.Response, error)
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
func (c *gitLabClient) Namespaces(opt *gitlab.ListNamespacesOptions, options ...gitlab.OptionFunc) ([]*gitlab.Namespace, *gitlab.Response, error) {
	namespaces, res, err := c.cli.Namespaces.ListNamespaces(opt, options...)
	return namespaces, res, err
}
