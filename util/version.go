package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func GetVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--dirty", "--always")
	gopath, set := os.LookupEnv("GOPATH")
	if !set {
		out, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			log.Fatal(err)
		}
		gopath = strings.TrimSuffix(string(out), "\n")
	}
	cmd.Dir = filepath.Join(gopath, "src", "github.com", "project-flogo", "cli")
	cmd.Env = append(os.Environ())

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("\\n")
	fc := re.ReplaceAllString(string(out), "")

	return fc
}
