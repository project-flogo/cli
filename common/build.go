package common

type BuildOptions struct {
	OptimizeImports bool
	EmbedConfig     bool
	BackupMain      bool
	BuildExist      bool
	Shim            string
}

type BuildPreProcessor interface {
	DoPreProcessing(project AppProject, options *BuildOptions) error
	Type() string
}

type BuildPostProcessor interface {
	DoPostProcessing(project AppProject, options *BuildOptions) error
}

var buildPreProcessors []BuildPreProcessor
var buildPostProcessors []BuildPostProcessor

func RegisterBuildPreProcessor(processor BuildPreProcessor) {
	buildPreProcessors = append(buildPreProcessors, processor)
}

func BuildPreProcessors() []BuildPreProcessor {
	return buildPreProcessors
}

func RegisterBuildPostProcessor(processor BuildPostProcessor) {
	buildPostProcessors = append(buildPostProcessors, processor)
}

func BuildPostProcessors() []BuildPostProcessor {
	return buildPostProcessors
}
