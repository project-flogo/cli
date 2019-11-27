package common

import "github.com/project-flogo/cli/util"

type AppProject interface {
	Validate() error
	Name() string
	Dir() string
	BinDir() string
	SrcDir() string
	Executable() string
	AddImports(ignoreError bool, addToJson bool, imports ...util.Import) error
	RemoveImports(imports ...string) error
	GetPath(flogoImport util.Import) (string, error)
	DepManager() util.DepManager

	GetGoImports(withVersion bool) ([]util.Import, error)
}
