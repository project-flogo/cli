package api

//Legacy Helper Functions
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/project-flogo/cli/common"
	"github.com/project-flogo/cli/util"
)

const (
	pkgLegacySupport = "github.com/project-flogo/legacybridge"
)

func InstallLegacySupport(project common.AppProject) error {
	//todo make sure we only install once
	pkgLegacySupportImport, err := util.NewFlogoImportFromPath(pkgLegacySupport)
	if err != nil {
		return err
	}
	err = project.AddImports(false, true, pkgLegacySupportImport)
	if err == nil {
		fmt.Println("Installed Legacy Support")
	}
	return err
}

func CreateLegacyMetadata(path, contribType, contribPkg string) error {

	var mdGoFilePath string

	tplMetadata := ""

	switch contribType {
	case "action":
		//ignore
		return nil
	case "trigger":
		fmt.Printf("Generating metadata for legacy trigger: %s\n", contribPkg)
		mdGoFilePath = filepath.Join(path, "trigger_metadata.go")
		tplMetadata = tplTriggerMetadataGoFile
	case "activity":
		fmt.Printf("Generating metadata for legacy actvity: %s\n", contribPkg)
		mdGoFilePath = filepath.Join(path, "activity_metadata.go")
		tplMetadata = tplActivityMetadataGoFile
	default:
		return nil
	}

	mdFilePath := filepath.Join(path, contribType+".json")
	pkg := filepath.Base(path)

	if idx := strings.Index(pkg, "@"); idx > 0 {
		pkg = pkg[:idx]
	}

	raw, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: string(raw),
	}

	err = os.Chmod(path, 0777)
	if err != nil {
		return err
	}
	defer os.Chmod(path, 0555)

	f, err := os.Create(mdGoFilePath)
	if err != nil {
		return err
	}
	RenderTemplate(f, tplMetadata, info)
	f.Close()

	return nil
}

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/project-flogo/legacybridge"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	legacybridge.RegisterLegacyActivity(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/project-flogo/legacybridge"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	legacybridge.RegisterLegacyTriggerFactory(md.ID, NewFactory(md))
}
`

//RenderTemplate renders the specified template
func RenderTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
