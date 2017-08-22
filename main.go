package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tomogoma/go-commons/config"
	"github.com/tomogoma/go-commons/errors"
	"github.com/tomogoma/repo-installer/dirs"
	"github.com/tomogoma/repo-installer/repositories"
)

func main() {

	var github = flag.String("github", "", "github.com account/name(s) separated by space e.g. \"tomogoma/micro-installer tomogoma/imagems\"")
	var bitbucket = flag.String("bitbucket", "", "bitbucket.org account/name(s) separated by space e.g. \"tomogoma/test tomogoma/test2\"")
	var file = flag.String("file", "", "/path/to/repositories.yml file containing github and bitbucket account/names as in repositiories_example.yml")
	var keepSrc = flag.Bool("keepSrc", false, "Provide this flag to keep src files on completion")
	var outputDir = flag.String("outputDir", "src", "path/to/parent-output dir to clone dir")
	var help = flag.Bool("help", false, "Print this help information")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	conf := &repositories.Config{Github: make([]string, 0), Bitbucket: make([]string, 0)}
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
	srcs := conf.Sources()
	if len(srcs) == 0 {
		flag.Usage()
		return
	}
	for _, source := range srcs {
		if err := fetchInstallSource(*outputDir, *keepSrc, source); err != nil {
			log.Println(err)
			return
		}
	}
}

func fetchInstallSource(outputDir string, keepSrc bool, source repositories.Source) error {
	baseURL, err := url.Parse(source.URL)
	if err != nil {
		return errors.Newf("Unable to parse '%s' as a URL: %v", source.URL, err)
	}
	repo := path.Join(baseURL.Path, source.Repo)
	dir := path.Join(outputDir, baseURL.Host, repo)
	if _, err := os.Stat(dir); err == nil {
		keepSrc = true
		err = pull(dir)
	} else {
		err = clone(baseURL, repo, dir)
	}
	defer cleanUp(keepSrc, dir)
	if err != nil {
		return errors.Newf("%v", err)
	}
	if err := install(dir); err != nil {
		return errors.Newf("error installing at '%s': %v", dir, err)
	}
	return nil
}

func pull(inDir string) error {
	dh := dirs.NewHelper()
	if err := dh.PushD(inDir); err != nil {
		return errors.Newf("error opening repo dir: %v", err)
	}
	defer func() {
		if err := dh.PopD(); err != nil {
			log.Printf("Error leaving repo dir: %v", err)
		}
	}()
	log.Printf("Pulling to update %s...", inDir)
	cmd := exec.Command("git", "pull")
	if err := cmd.Run(); err != nil {
		return errors.Newf("unable to pull: %v", err)
	}
	log.Printf("Done pulling.")
	return nil
}

func clone(baseURL *url.URL, repo string, intoDir string) error {
	if err := os.MkdirAll(intoDir, 0755); err != nil {
		return errors.Newf("error creating source dir: %v", err)
	}
	URL := baseURL.ResolveReference(&url.URL{Path: repo})
	log.Printf("Cloning %s into %s...", URL, intoDir)
	cmd := exec.Command("git", "clone", URL.String(), intoDir)
	if err := cmd.Run(); err != nil {
		return errors.Newf("unable to clone: %v", repo, err)
	}
	log.Printf("Done cloning.")
	return nil
}

func install(dir string) error {
	log.Printf("Installing %s...", dir)
	dh := dirs.NewHelper()
	if err := dh.PushD(dir); err != nil {
		return errors.Newf("error changing to installer dir: %v", err)
	}
	defer func() {
		if err := dh.PopD(); err != nil {
			log.Printf("error switching back directory: %v", err)
		}
	}()
	cmd := exec.Command("make", "install")
	err := cmd.Run()
	if err != nil {
		return errors.Newf("error running installer script: %v", err)
	}
	log.Printf("Done Installing.")
	return nil
}

func cleanUp(keepSrc bool, dir string) {
	if !keepSrc {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("Error cleaning up: %v", err)
		}
	}
}
