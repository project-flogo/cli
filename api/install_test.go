package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var newJsonString = `{
	"name": "temp",
	"type": "flogo:app",
	"version": "0.0.1",
	"description": "My flogo application description",
	"appModel": "1.0.0",
	"imports": [
	  "github.com/project-flogo/flow",
	  "github.com/skothari-tibco/flogoaztrigger",
	  "github.com/project-flogo/contrib/activity/actreturn",
	  "github.com/project-flogo/contrib/activity/log",
	  "github.com/TIBCOSoftware/flogo-contrib/activity/rest"
	],
	"triggers": [
	  {
		"id": "my_rest_trigger",
		"ref":  "github.com/skothari-tibco/flogoaztrigger",
		"handlers": [
		  {
			"action": {
			  "ref": "github.com/project-flogo/flow",
			  "settings": {
				"flowURI": "res://flow:simple_flow"
			  },
			  "input": {
				"in": "inputA"
			  },
			  "output" :{
						"out": "=$.out"
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
				  "message": "=$flow.in",
				  "flowInfo": "false",
				  "addToFlow": "false"
				}
			  }
			},
			{
				"id" :"return",
				"name" : "Activity Return",
				"activity":{
					"ref" : "github.com/project-flogo/contrib/activity/actreturn",
					"settings":{
						"mappings":{
							"out": "nameA"
						}
					}
				}
			}
		  ],
		  "links": [
			  {
				  "from":"log",
				  "to":"return"
			  }
		  ]
		}
	  }
	]
  }
  `

func TestInstallLegacyPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/TIBCOSoftware/flogo-contrib/activity/log")
	assert.Equal(t, nil, err)

}

func TestInstallPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/skothari-tibco/csvtimer")
	assert.Equal(t, nil, err)

}
func TestInstallPkgWithVersion(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/project-flogo/contrib/activity/log@v0.9.0-alpha.3")
	assert.Equal(t, nil, err)

}
func TestListPkg(t *testing.T) {
	t.Log("Testing listing of packages")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListPackages(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, false)
	assert.Equal(t, nil, err)

}

func TestListAllPkg(t *testing.T) {
	t.Log("Testing listing of all contribs")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListPackages(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, true)
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

	err = ListPackages(NewAppProject(filepath.Join(testEnv.currentDir, "temp")), true, false)
	assert.Equal(t, nil, err)
}
