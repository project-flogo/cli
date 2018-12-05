package api

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallLegacyPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/TIBCOSoftware/flogo-contrib/activity/log")
	assert.Equal(t, nil, err)

}

func TestInstallPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/skothari-tibco/csvtimer")
	assert.Equal(t, nil, err)

}

func TestListPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListPackages(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, false)
	assert.Equal(t, nil, err)

}

func TestListAllPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListPackages(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, true)
	assert.Equal(t, nil, err)

}
