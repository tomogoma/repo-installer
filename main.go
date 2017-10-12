package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"

	"github.com/tomogoma/go-commons/config"
	"github.com/tomogoma/go-commons/errors"
	"github.com/tomogoma/repo-installer/dirs"
	"github.com/tomogoma/repo-installer/repositories"
)

func main() {

	var file = flag.String("file", "", "/path/to/repositories.yml file containing repository urls.")
	var keepSrc = flag.Bool("keepSrc", false, "Provide this flag to keep src files on completion")
	var outputDir = flag.String("outputDir", "src", "path/to/parent-output dir to clone dir")
	var help = flag.Bool("help", false, "Print this help information")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	conf := repositories.Config{}
	if err := config.ReadYamlConfig(*file, &conf); err != nil {
		log.Printf("error reading repositories file: %v", err)
		flag.Usage()
		return
	}

	conf.Clean()
	if len(conf.Repos) == 0 {
		flag.Usage()
		return
	}

	for _, repo := range conf.Repos {
		if err := run(*outputDir, *keepSrc, repo); err != nil {
			log.Println(err)
		}
	}
}

func run(outputDir string, keepSrc bool, source string) error {
	URL, err := url.Parse(source)
	if err != nil {
		return errors.Newf("Unable to parse '%s' as a URL: %v", source, err)
	}
	dir := path.Join(outputDir, URL.Host, URL.Path)
	defer cleanUp(keepSrc, dir)
	if _, err := os.Stat(dir); err == nil {
		keepSrc = true
		err = pull(dir)
	} else {
		err = clone(dir, source)
	}
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
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("unable to pull: %v: %s", err, out)
	}
	log.Printf("Done pulling.")
	return nil
}

func clone(intoDir, repo string) error {
	if err := os.MkdirAll(intoDir, 0755); err != nil {
		return errors.Newf("error creating source dir: %v", err)
	}
	log.Printf("Cloning %s into %s...", repo, intoDir)
	cmd := exec.Command("git", "clone", repo, intoDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("unable to clone: %v: %s", err, out)
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
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("error running installer script: %v: %s", err, out)
	}
	log.Printf("Done Installing.")
	return nil
}

func cleanUp(keepSrc bool, dir string) {
	if keepSrc {
		return
	}
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("Error cleaning up: %v", err)
	}
}
