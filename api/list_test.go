package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllContribs(t *testing.T) {
	t.Log("Testing listing of all contribs")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	os.Chdir(testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListContribs(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, "all")
	assert.Equal(t, nil, err)

}

func TestListWithLegacyPkg(t *testing.T) {
	t.Log("Testing listing of legacy contribs")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	err := os.Chdir(tempDir)

	file, err := os.Create("flogo.json")
	if err != nil {
		t.Fatal(err)
		assert.Equal(t, true, false)
	}
	defer file.Close()
	fmt.Fprintf(file, newJsonString)
	_, err = CreateProject(testEnv.currentDir, "temp", "flogo.json", "")
	assert.Equal(t, nil, err)

	err = ListContribs(NewAppProject(filepath.Join(testEnv.currentDir, "temp")), true, "")
	assert.Equal(t, nil, err)
}
