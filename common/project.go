package common

type AppProject interface {
	Validate() error
	Name() string
	Dir() string
	BinDir() string
	SrcDir() string
	Executable() string
	AddImports(imports ...string) error
	RemoveImports(imports ...string) error
	GetPath(pkg string) (string, error)
}
