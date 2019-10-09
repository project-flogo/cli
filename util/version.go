package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
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

func CreateVersionFile(cmdPath, currentVersion string) error {

	f, err := os.Create(filepath.Join(cmdPath, "currentversion.go"))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer f.Close()

	_ = packageTemplate.Execute(f, struct {
		Timestamp time.Time
		Version   string
	}{
		Timestamp: time.Now(),
		Version:   currentVersion,
	})

	return nil
}

var packageTemplate = template.Must(template.New("").Parse(`// Generated Code; DO NOT EDIT.
// {{ .Timestamp }}
package main

func init() {
	Version = "{{ .Version }}"
}
`))
