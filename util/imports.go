package util

import (
	"errors"
	"fmt"
	"path"
	"regexp"
)

/* util.Import struct defines the different fields which can be extracted from a Flogo import
these imports are stored in flogo.json in the "imports" array, for instance:

 "imports": [
   "github.com/project-flogo/contrib@v0.9.0-alpha.4:/activity/log",
   "github.com/project-flogo/contrib/activity/rest@v0.9.0"
   "rest_activity github.com/project-flogo/contrib@v0.9.0:/activity/rest",
   "rest_trigger github.com/project-flogo/contrib:/trigger/rest",
   "github.com/project-flogo/flow"
 ]

*/

type FlogoImport struct {
	modulePath         string
	relativeImportPath string
	version            string
	alias              string
}

func NewFlogoImportFromPath(flogoImportPath string) (Import, error) {
	flogoImport, err := ParseImport(flogoImportPath)
	if err != nil {
		return nil, err
	}
	return flogoImport, nil
}

func NewFlogoImport(modulePath, relativeImportPath, version, alias string) Import {
	return &FlogoImport{modulePath: modulePath, relativeImportPath: relativeImportPath, version: version, alias: alias}
}

func NewFlogoImportWithVersion(imp Import, version string) Import {
	return &FlogoImport{modulePath: imp.ModulePath(), relativeImportPath: imp.RelativeImportPath(), version: version, alias: imp.Alias()}
}

type Import interface {
	fmt.Stringer

	ModulePath() string
	RelativeImportPath() string
	Version() string
	Alias() string

	CanonicalImport() string // canonical import is used in flogo.json imports array and to check for equality of imports
	GoImportPath() string    // the import path used in .go files
	GoGetImportPath() string // the import path used by "go get" command
	GoModImportPath() string // the import path used by "go mod edit" command
	IsClassic() bool         // an import is "classic" if it has no : character separator, hence no relative import path
	CanonicalAlias() string  // canonical alias is the alias used in the flogo.json
}

type Imports []Import

func (flogoImport *FlogoImport) ModulePath() string {
	return flogoImport.modulePath
}

func (flogoImport *FlogoImport) RelativeImportPath() string {
	return flogoImport.relativeImportPath
}

func (flogoImport *FlogoImport) Version() string {
	return flogoImport.version
}

func (flogoImport *FlogoImport) Alias() string {
	return flogoImport.alias
}

func (flogoImport *FlogoImport) CanonicalImport() string {
	alias := ""
	if flogoImport.alias != "" {
		alias = flogoImport.alias + " "
	}
	version := ""
	if flogoImport.version != "" {
		version = "@" + flogoImport.version
	}
	relativeImportPath := ""
	if flogoImport.relativeImportPath != "" {
		relativeImportPath = ":" + flogoImport.relativeImportPath
	}

	return alias + flogoImport.modulePath + version + relativeImportPath
}

func (flogoImport *FlogoImport) CanonicalAlias() string {
	if flogoImport.alias != "" {
		return flogoImport.alias
	} else {
		return path.Base(flogoImport.GoImportPath())
	}
}

func (flogoImport *FlogoImport) GoImportPath() string {
	return flogoImport.modulePath + flogoImport.relativeImportPath
}

func (flogoImport *FlogoImport) GoGetImportPath() string {
	version := "@latest"
	if flogoImport.version != "" {
		version = "@" + flogoImport.version
	}
	return flogoImport.modulePath + flogoImport.relativeImportPath + version
}

func (flogoImport *FlogoImport) GoModImportPath() string {
	version := "@latest"
	if flogoImport.version != "" {
		version = "@" + flogoImport.version
	}
	return flogoImport.modulePath + version
}

func (flogoImport *FlogoImport) IsClassic() bool {
	return flogoImport.relativeImportPath == ""
}

func (flogoImport *FlogoImport) String() string {
	version := ""
	if flogoImport.version != "" {
		version = " " + flogoImport.version
	}
	relativeImportPath := ""
	if flogoImport.relativeImportPath != "" {
		relativeImportPath = flogoImport.relativeImportPath
	}

	return flogoImport.modulePath + relativeImportPath + version
}

var flogoImportPattern = regexp.MustCompile(`^(([^ ]*)[ ]+)?([^@:]*)@?([^:]*)?:?(.*)?$`) // extract import path even if there is an alias and/or a version

func ParseImports(flogoImports []string) (Imports, error) {
	var result Imports

	for _, flogoImportPath := range flogoImports {
		flogoImport, err := ParseImport(flogoImportPath)
		if err != nil {
			return nil, err
		}
		result = append(result, flogoImport)
	}

	return result, nil
}

func ParseImport(flogoImport string) (Import, error) {
	if !flogoImportPattern.MatchString(flogoImport) {
		return nil, errors.New(fmt.Sprintf("The Flogo import '%s' cannot be parsed.", flogoImport))
	}

	matches := flogoImportPattern.FindStringSubmatch(flogoImport)

	result := &FlogoImport{modulePath: matches[3], relativeImportPath: matches[5], version: matches[4], alias: matches[2]}

	return result, nil
}
