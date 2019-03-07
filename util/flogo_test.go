package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	line := `              "ref":"github.com/TIBCOSoftware/flogo-contrib/activity/log",`

	if idx := strings.Index(line, "\"ref\""); idx > -1 {

		fmt.Printf("line: '%s'\n", line)

		startPkgIdx := strings.Index(line[idx+6:], "\"")

		fmt.Printf("line frag: '%s'\n", line[idx+6+startPkgIdx:])
		pkg := strings.Split(line[idx+6+startPkgIdx:], "\"")[1]
		//pkg := strings.Split(",")[0]
		//pkg := strings.Split(line, ":")[1]
		//pkg = strings.TrimSpace(pkg)
		//pkg = pkg[1 : len(pkg)-2]
		fmt.Printf("package: '%s'\n", pkg)
	}
}
