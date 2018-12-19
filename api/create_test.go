package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonString = `{
    "name": "_APP_NAME_",
    "type": "flogo:app",
    "version": "0.0.1",
    "description": "My flogo application description",
    "appModel": "1.0.0",
    "triggers": [
      {
        "id": "my_rest_trigger",
        "ref": "github.com/project-flogo/contrib/trigger/rest",
        "settings": {
          "port": "8888"
        },
        "handlers": [
          {
            "settings": {
              "method": "GET",
              "path": "/test/:val"
            },
            "action": {
              "ref": "github.com/project-flogo/flow",
              "settings": {
                "flowURI": "res://flow:simple_flow"
              },
              "input": {
                "in": "$.pathParams.val"
              }
            }
          }
        ]
      }
    ],
    "resources": [
      {
        "id": "flow:simple_flow",
        "data": {
          "name": "simple_flow",
          "metadata": {
            "input": [
              { "name": "in", "type": "string",  "value": "test" }
            ],
            "output": [
              { "name": "out", "type": "string" }
            ]
          },
          "tasks": [
            {
              "id": "log",
              "name": "Log Message",
              "activity": {
                "ref": "github.com/project-flogo/contrib/activity/log",
                "input": {
                  "message": "$flow.in",
                  "flowInfo": "false",
                  "addToFlow": "false"
                }
              }
            }
          ],
          "links": []
        }
      }
    ]
  }
  `

type TestEnv struct {
	currentDir string
}

func (t *TestEnv) getTestwd() (dir string, err error) {
	return t.currentDir, nil
}

func (t *TestEnv) cleanup() {

	os.RemoveAll(t.currentDir)
}

func GetTempDir() string {
	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("GO111MODULE", "on")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Unsetenv("FLOGO_BUILD_EXPERIMENTAL")

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	tempDirInfo, err := filepath.EvalSymlinks(tempDir)
	if err == nil {
		// Sym link
		tempDir = tempDirInfo
	}
	return tempDir
}
func TestCmdCreate_noflag(t *testing.T) {
	t.Log("Testing simple creation of project")

	tempDir := GetTempDir()
	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	assert.Equal(t, nil, CreateProject(testEnv.currentDir, "myApp", "", ""))

	_, err := os.Stat(filepath.Join(tempDir, "myApp", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "myApp", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "main.go"))
	assert.Equal(t, nil, err)

}

func TestCmdCreate_flag(t *testing.T) {
	t.Log("Testing creation of project while the file is provided")

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
	os.Chdir(testEnv.currentDir)
	file, err := os.Create("flogo.json")
	if err != nil {
		t.Fatal(err)
		assert.Equal(t, true, false)
	}
	defer file.Close()
	fmt.Fprintf(file, jsonString)
	assert.Equal(t, nil, CreateProject(testEnv.currentDir, "flogo", "flogo.json", ""))

	_, err = os.Stat(filepath.Join(tempDir, "flogo", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "flogo", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "flogo", "src", "main.go"))
	assert.Equal(t, nil, err)

}

func TestCmdCreate_masterCore(t *testing.T) {
	t.Log("Testing creation of project when the version of core is provided `master`")

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

	assert.Equal(t, nil, CreateProject(testEnv.currentDir, "myApp", "", "master"))

}

func TestCmdCreate_versionCore(t *testing.T) {
	t.Log("Testing creation of project when the version of core is provided `v0.9.0-alpha.3`")
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

	assert.Equal(t, nil, CreateProject(testEnv.currentDir, "myApp", "", "v0.9.0-alpha.3"))

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "myApp", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "main.go"))
	assert.Equal(t, nil, err)

	data, err1 := ioutil.ReadFile(filepath.Join(tempDir, "myApp", "src", "go.mod"))
	assert.Equal(t, nil, err1)

	assert.Equal(t, true, strings.Contains(string(data), "v0.9.0-alpha.3"))
}
