package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pborman/uuid"
	"github.com/tomogoma/go-commons/config"
	"github.com/tomogoma/go-commons/errors"
)

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

var github = flag.String("github", "", "github.com account/name(s) separated by space e.g. \"tomogoma/micro-installer tomogoma/imagems\"")
var bitbucket = flag.String("bitbucket", "", "bitbucket.org account/name(s) separated by space e.g. \"tomogoma/test tomogoma/test2\"")
var file = flag.String("file", "", "/path/to/repositories.yml file containing github and bitbucket account/names as in repositiories_example.yml")
var keepSrc = flag.Bool("keepSrc", false, "Provide this flag to keep src files on completion")
var outputDir = flag.String("outputDir", "src", "path/to/parent-output dir to clone dir")

func init() {
	flag.Parse()
}

func main() {
	conf := &Config{Github: make([]string, 0), Bitbucket: make([]string, 0)}
	if *file != "" {
		if err := config.ReadYamlConfig(*file, conf); err != nil {
			log.Printf("error reading repositories file: %v", err)
			return
		}
	}
	if *github != "" {
		ghs := strings.Split(*github, " ")
		conf.Github = append(conf.Github, ghs...)
	}
	if *bitbucket != "" {
		bbs := strings.Split(*bitbucket, " ")
		conf.Bitbucket = append(conf.Bitbucket, bbs...)
	}
	conf.clean()
	if len(conf.Github) == 0 && len(conf.Bitbucket) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Unable to establish location of executable: %v", err)
		return
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		log.Printf("Unable to establish location of executable (while resolving symlinks): %v", err)
		return
	}
	UUID := uuid.New()
	dirs, err := clone(*outputDir, conf.Github, "https://github.com", UUID)
	if err != nil {
		log.Printf("Error cloning github: %v", err)
		return
	}
	defer func() {
		cleanUp(*keepSrc, dirs)
	}()
	bDirs, err := clone(*outputDir, conf.Bitbucket, "https://bitbucket.org", UUID)
	if err != nil {
		log.Printf("Error cloning bitbucket: %v", err)
		return
	}
	dirs = append(dirs, bDirs...)
	for _, dir := range dirs {
		if err := install(dir); err != nil {
			log.Printf("error installing at '%s': %v", dir, err)
		}
	}
}

func clone(srcDir string, repos []string, baseURL string, UUID string) ([]string, error) {
	bURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.Newf("Unable to parse '%s' as a URL: %v", baseURL, err)
	}
	dirs := make([]string, 0)
	for _, repo := range repos {
		dir := path.Join(srcDir, bURL.Host, bURL.Path, repo+"_"+UUID)
		repo = path.Join(bURL.Path, repo)
		URL := bURL.ResolveReference(&url.URL{Path: repo})
		if _, err := os.Stat(dir); err == nil {
			log.Printf("Skipping '%s' - already exists", repo)
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, errors.Newf("error creating source dir: %v", err)
		}
		log.Printf("Cloning %s into %s...", URL, dir)
		cmd := exec.Command("git", "clone", URL.String(), dir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, errors.Newf("unable to clone: %v", repo, err)
		}
		log.Printf("%s", out)
		dirs = append(dirs, dir)
		log.Printf("Done cloning.")
	}
	return dirs, nil
}

func install(dir string) error {
	log.Printf("Installing %s...", dir)
	execF, err := os.Executable()
	if err != nil {
		return errors.Newf("error getting executable location: %v", err)
	}
	workingDir := filepath.Dir(execF)
	if err := os.Chdir(dir); err != nil {
		return errors.Newf("error switching to install dir: %v", err)
	}
	defer func() {
		os.Chdir(workingDir)
	}()
	cmd := exec.Command("make", "install")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Newf("error running installer script: %v", err)
	}
	log.Printf("%s", out)
	log.Printf("Done Installing.")
	return nil
}

func cleanUp(keepSrc bool, dirs []string) {
	for _, dir := range dirs {
		if keepSrc {
			i := strings.LastIndex(dir, "_")
			newDir := dir[0:i]
			if err := os.Rename(dir, newDir); err != nil {
				log.Printf("Error renaming repo (%s) to natural name (...%s): %v",
					dir, path.Base(newDir), err)
			}
		} else {
			if err := os.RemoveAll(dir); err != nil {
				log.Printf("Error removing repo '%s': %v", dir, err)
			}
		}
	}
}
