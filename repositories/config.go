package repositories

import "strings"

type Source struct {
	URL  string
	Repo string
}

type Config struct {
	Github    []string `yaml:"github"`
	Bitbucket []string `yaml:"bitbucket"`
}

func (c *Config) clean() {
	if c.Github == nil {
		c.Github = make([]string, 0)
	}
	if c.Bitbucket == nil {
		c.Bitbucket = make([]string, 0)
	}
	github := make(map[string]int)
	bitbucket := make(map[string]int)
	for i, path := range c.Github {
		if strings.TrimSpace(path) == "" {
			continue
		}
		github[path] = i
	}
	for i, path := range c.Bitbucket {
		if strings.TrimSpace(path) == "" {
			continue
		}
		bitbucket[path] = i
	}
	c.Github = make([]string, 0)
	c.Bitbucket = make([]string, 0)
	for path := range github {
		c.Github = append(c.Github, path)
	}
	for path := range bitbucket {
		c.Bitbucket = append(c.Bitbucket, path)
	}
}

func (c *Config) Sources() []Source {
	c.clean()
	ss := make([]Source, 0)
	for _, repo := range c.Github {
		ss = append(ss, Source{URL: "https://github.com", Repo: repo})
	}
	for _, repo := range c.Bitbucket {
		ss = append(ss, Source{URL: "https://bitbucket.org", Repo: repo})
	}
	return ss
}
