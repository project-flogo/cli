package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"
)

var goPathCached string

func FindOldPackageSrc(pkg string) (srcPath, srcVer string, err error) {

	goPath := GetGoPath()

	pkgParts := strings.Split(pkg, "/")
	path := filepath.Join(goPath, "src", filepath.Join(pkgParts...))

	if _, e := os.Stat(path); !os.IsNotExist(e) {
		// path/to/whatever exists
		v := GetPackageVersionFromSource(pkg)

		return path, v, nil
	}

	return "", "", newPkgNotFoundError(pkg)
}

func FindGoModPackageSrc(pkg string, version string, latest bool) (srcPath, srcVer string, err error) {

	pkgParts := strings.Split(pkg, "/")
	if len(pkgParts) < 2 {
		return "", "", fmt.Errorf("invalid package: %s", pkg)
	}

	name := pkgParts[len(pkgParts)-1]
	path := pkgParts[:len(pkgParts)-1]

	goPath := GetGoPath()
	flogoPkgPath := filepath.Join(goPath, "pkg", "mod", filepath.Join(path...))

	if _, e := os.Stat(flogoPkgPath); os.IsNotExist(e) {
		return "", "", newPkgNotFoundError(pkg)
	}

	var files []string

	err = filepath.Walk(flogoPkgPath, visit(name, &files))
	if err != nil {
		return "", "", newPkgNotFoundError(pkg)
	}

	if latest {

		lf := ""
		var lv *semver.Version

		for _, file := range files {

			parts := strings.SplitN(file, "@v", 2)

			if lf == "" {
				lf = file
				lv, err = semver.NewVersion(parts[1])
				if err != nil {
					return "", "", err
				}
				continue
			}

			sv, err := semver.NewVersion(parts[1])
			if err != nil {
				return "", "", err
			}

			if lv.LessThan(*sv) {
				lf = file
				lv = sv
			}
		}

		if lf == "" {
			return "", "", newPkgNotFoundError(pkg)
		}

		return lf, lv.String(), nil
	} else {

		for _, file := range files {

			parts := strings.SplitN(file, "@v", 2)
			if parts[1] == version {
				return file, version, nil
			}
		}
	}

	return "", "", newPkgNotFoundError(pkg)
}

func visit(name string, files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			return nil
		}

		if strings.HasPrefix(info.Name(), name+"@v") {
			*files = append(*files, path)
		}

		return nil
	}
}

func GetGoPath() string {

	if goPathCached != "" {
		return goPathCached
	}

	set := false
	goPathCached, set = os.LookupEnv("GOPATH")
	if !set {
		out, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			log.Fatal(err)
		}
		goPathCached = strings.TrimSuffix(string(out), "\n")
	}

	return goPathCached
}

func newPkgNotFoundError(pkg string) error {
	return &pkgNotFoundError{pkg: pkg}
}

type pkgNotFoundError struct {
	pkg string
}

func (e *pkgNotFoundError) Error() string {
	return fmt.Sprintf("Package '%s' not found", e.pkg)
}

func IsPkgNotFoundError(err error) bool {
	_, ok := err.(*pkgNotFoundError)
	return ok
}