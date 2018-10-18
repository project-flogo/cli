package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEnv struct {
	currentDir string
}

func (t *TestEnv) getTestwd() (dir string, err error) {
	return t.currentDir, nil
}

func (t *TestEnv) cleanup() {

	os.RemoveAll(t.currentDir)
}

func TestCmdCreate_noflag(t *testing.T) {

	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
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

	assert.Equal(t, nil, CreateProject("myApp", false, "", tempDir))

	_, err = os.Stat(Concat(tempDir, "/myApp/src/go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(Concat(tempDir, "/myApp/flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(Concat(tempDir, "/myApp/src/main.go"))
	assert.Equal(t, nil, err)

	testEnv.cleanup()
}

func TestCmdCreate_flag(t *testing.T) {

	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
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

	cliCmd, err := exec.Command("cp", "/Users/skothari-tibco/Desktop/flogo.json", tempDir).CombinedOutput()
	if err != nil {
		fmt.Println(string(cliCmd))
		assert.Equal(t, true, false)
	}
	assert.Equal(t, nil, CreateProject("flogo.json", true, "", tempDir))

	_, err = os.Stat(Concat(tempDir, "/flogo/src/go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(Concat(tempDir, "/flogo/flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(Concat(tempDir, "/flogo/src/main.go"))
	assert.Equal(t, nil, err)

}

func TestCmdCreate_masterCore(t *testing.T) {

	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
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

	assert.Equal(t, nil, CreateProject("myApp", false, "master", tempDir))

}

func TestCmdCreate_versionCore(t *testing.T) {

	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
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

	assert.Equal(t, nil, CreateProject("myApp", false, "v0.9.0-alpha.0", tempDir))

	_, err = os.Stat(Concat(tempDir, "/myApp/src/go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(Concat(tempDir, "/myApp/flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(Concat(tempDir, "/myApp/src/main.go"))
	assert.Equal(t, nil, err)

	data, err1 := ioutil.ReadFile(Concat(tempDir, "/myApp/src/go.mod"))
	assert.Equal(t, nil, err1)

	assert.Equal(t, strings.Contains(string(data), "v0.9.0-alpha.0"), true)
}
