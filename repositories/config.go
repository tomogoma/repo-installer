package repositories

import "strings"

type Config struct {
	Repos []string `yaml:"repos"`
}

func (c Config) Clean() {
	if len(c.Repos) == 0 {
		return
	}
	var repos []string
	for _, repo := range c.Repos {
		if strings.TrimSpace(repo) == "" {
			continue
		}
		repos = append(repos, repo)
	}
	c.Repos = repos
}
