package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	cliPackage = "github.com/project-flogo/cli"
)

func GetCLIInfo() (string, string, error) {

	path, ver, err := FindOldPackageSrc(cliPackage)

	if IsPkgNotFoundError(err) {
		//must be using the new go mod layout
		path, ver, err = FindGoModPackageSrc(cliPackage, "", true)
	}

	return path, ver, err
}

func GetPackageVersionOld(pkg string) string {
	re := regexp.MustCompile("\\n")

	cmd := exec.Command("git", "describe", "--tags", "--dirty", "--always")
	cmd.Env = append(os.Environ())

	gopath := GetGoPath()

	pkgParts := strings.Split(pkg, "/")
	cmd.Dir = filepath.Join(gopath, "src", filepath.Join(pkgParts...))

	out, err := cmd.Output() // execute "git describe"
	if err != nil {
		log.Fatal(err)
	}
	fc := re.ReplaceAllString(string(out), "")

	if len(fc) > 1 {
		return fc[1:]
	}

	return fc
}
