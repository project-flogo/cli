// +build ignore

package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"github.com/project-flogo/cli/util"
)

// This Go program is aimed at being called by go:generate from "cmd/flogo/main.go" to create a "currentversion.go" file
// in the "cmd/flogo" subdirectory.
//
// Once this file is created, it will set at runtime the version of the CLI without relying on GOPATH nor Git command to
// parse the version from the source. Hence it is possible to distribute the CLI as a fully static binary.
//
// Users getting the CLI with a classic "go get" command will still have the version retrieved from the directory
// $GOPATH/src/github.com/project-flogo/cli
func main() {
	currentVersion := util.GetVersion(false)

	f, err := os.Create("./currentversion.go")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	packageTemplate.Execute(f, struct {
		Timestamp time.Time
		Version   string
	}{
		Timestamp: time.Now(),
		Version:   currentVersion,
	})
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// {{ .Timestamp }}
package main

func init() {
	Version = "{{ .Version }}"
}
`))
