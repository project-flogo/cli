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
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Env = append(os.Environ())

	re := regexp.MustCompile("\\n")

	currentRemoteURLOutput, err := cmd.Output() // determine whether we're building from source
	currentRemoteURL := re.ReplaceAllString(string(currentRemoteURLOutput), "")

	cmd = exec.Command("git", "describe", "--tags", "--dirty", "--always")

	if !strings.HasSuffix(currentRemoteURL, "cli.git") { // we're not building from source but we are "go getting"
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
