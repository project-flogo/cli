package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallPkg(t *testing.T) {
	t.Log("Testing installation of package")

	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO111MODULE", "on")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("FLOGO_BUILD_EXPERIMENTAL")

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	tempDirInfo, err := filepath.EvalSymlinks(tempDir)
	if err == nil {
		// Sym link
		tempDir = tempDirInfo
	}

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err = CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "myApp", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "main.go"))
	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(testEnv.currentDir), "github.com/TIBCOSoftware/flogo-contrib/activity/log")
	assert.Equal(t, nil, err)

}
