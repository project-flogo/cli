package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func GetVersion(fromGoPathSources bool) string {
	re := regexp.MustCompile("\\n")

	cmd := exec.Command("git", "describe", "--tags", "--dirty", "--always")
	cmd.Env = append(os.Environ())

	if fromGoPathSources {
		gopath, set := os.LookupEnv("GOPATH")
		if !set {
			out, err := exec.Command("go", "env", "GOPATH").Output()
			if err != nil {
				log.Fatal(err)
			}
			gopath = strings.TrimSuffix(string(out), "\n")
		}
		cmd.Dir = filepath.Join(gopath, "src", "github.com", "project-flogo", "cli")
	}

	out, err := cmd.Output() // execute "git describe"
	if err != nil {
		log.Fatal(err)
	}
	fc := re.ReplaceAllString(string(out), "")

	return fc
}
