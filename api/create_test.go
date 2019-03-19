package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonString = `{
  "name": "_APP_NAME_",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "My flogo application description",
  "appModel": "1.0.0",
  "imports": [
    "github.com/TIBCOSoftware/flogo-contrib/activity/log",
    "github.com/project-flogo/contrib/trigger/rest",
    "github.com/project-flogo/flow"
  ],
  "triggers": [
    {
      "id": "my_rest_trigger",
      "type": "rest",
      "settings": {
        "port": "8888"
      },
      "handlers": [
        {
          "settings": {
            "method": "GET",
            "path": "/test/:val"
          },
          "actions": [
            {
              "type": "flow",
              "settings": {
                "flowURI": "res://flow:simple_flow"
              },
              "input": {
                "in": "=$.pathParams.val"
              }
	    			}	
	  			]
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
              "type": "log",
              "input": {
                "message": "=$flow.in",
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

func TestCmdCreate_noflag(t *testing.T) {
	t.Log("Testing simple creation of project")

	tempDir, _ := GetTempDir()
	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")
	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "myApp", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "main.go"))
	assert.Equal(t, nil, err)
}

func TestCmdCreate_flag(t *testing.T) {
	t.Log("Testing creation of project while the file is provided")

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
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
	_, err = CreateProject(testEnv.currentDir, "flogo", "flogo.json", "")
	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "flogo", "src", "go.mod"))

	assert.Equal(t, nil, err)
	_, err = os.Stat(filepath.Join(tempDir, "flogo", "flogo.json"))

	assert.Equal(t, nil, err)

	_, err = os.Stat(filepath.Join(tempDir, "flogo", "src", "main.go"))
	assert.Equal(t, nil, err)
}

func TestCmdCreate_masterCore(t *testing.T) {
	t.Log("Testing creation of project when the version of core is provided `master`")

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	os.Chdir(testEnv.currentDir)

	_, err = CreateProject(testEnv.currentDir, "myApp", "", "master")
	assert.Equal(t, nil, err)
}

//todo fix this test, unreliable
//func TestCmdCreate_versionCore(t *testing.T) {
//	t.Log("Testing creation of project when the version of core is provided `v0.9.0-alpha.4`")
//
//	tempDir, err := ioutil.TempDir("", "test")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	testEnv := &TestEnv{currentDir: tempDir}
//
//	defer testEnv.cleanup()
//
//	t.Logf("Current dir '%s'", testEnv.currentDir)
//	os.Chdir(testEnv.currentDir)
//
//	_, err = CreateProject(testEnv.currentDir, "myApp", "", "v0.9.0-alpha.4")
//	assert.Equal(t, nil, err)
//
//	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "go.mod"))
//
//	assert.Equal(t, nil, err)
//	_, err = os.Stat(filepath.Join(tempDir, "myApp", "flogo.json"))
//
//	assert.Equal(t, nil, err)
//
//	_, err = os.Stat(filepath.Join(tempDir, "myApp", "src", "main.go"))
//	assert.Equal(t, nil, err)
//
//	data, err1 := ioutil.ReadFile(filepath.Join(tempDir, "myApp", "src", "go.mod"))
//	assert.Equal(t, nil, err1)
//
//	//todo fix, not a reliable test giving that importing latest of flow which will affect this import
//	assert.Equal(t, true, strings.Contains(string(data), "v0.9.0-alpha.4"))
//	fmt.Println(string(data))
//
//	appProject := NewAppProject(filepath.Join(testEnv.currentDir, "myApp"))
//
//	err = appProject.Validate()
//	assert.Nil(t, err)
//
//	common.SetCurrentProject(appProject)
//
//	err = BuildProject(common.CurrentProject(), BuildOptions{})
//	assert.Nil(t, err)
//}
