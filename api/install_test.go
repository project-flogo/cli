package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/project-flogo/cli/common"
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
	_ = os.Chdir(testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "v0.9.2")

	assert.Nil(t, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/TIBCOSoftware/flogo-contrib/activity/log")
	assert.Nil(t, err)

	appProject := NewAppProject(filepath.Join(testEnv.currentDir, "myApp"))

	err = appProject.Validate()
	assert.Nil(t, err)

	common.SetCurrentProject(appProject)

	err = BuildProject(common.CurrentProject(), common.BuildOptions{})
	assert.Nil(t, err)

}

func TestInstallPkg(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	_ = os.Chdir(testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Nil(t, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/project-flogo/contrib/activity/noop")
	assert.Nil(t, err)

	appProject := NewAppProject(filepath.Join(testEnv.currentDir, "myApp"))

	err = appProject.Validate()
	assert.Nil(t, err)

	common.SetCurrentProject(appProject)

	err = BuildProject(common.CurrentProject(), common.BuildOptions{})
	assert.Nil(t, err)
}

func TestInstallPkgWithVersion(t *testing.T) {
	t.Log("Testing installation of package")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	_ = os.Chdir(testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Nil(t, err)

	err = InstallPackage(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), "github.com/project-flogo/contrib/activity/log@v0.9.0")
	assert.Nil(t, err)

	appProject := NewAppProject(filepath.Join(testEnv.currentDir, "myApp"))

	err = appProject.Validate()
	assert.Nil(t, err)

	common.SetCurrentProject(appProject)

	err = BuildProject(common.CurrentProject(), common.BuildOptions{})
	assert.Nil(t, err)
}

func TestListPkg(t *testing.T) {
	t.Log("Testing listing of packages")

	tempDir, _ := GetTempDir()

	testEnv := &TestEnv{currentDir: tempDir}

	defer testEnv.cleanup()

	t.Logf("Current dir '%s'", testEnv.currentDir)
	_ = os.Chdir(testEnv.currentDir)

	_, err := CreateProject(testEnv.currentDir, "myApp", "", "")

	assert.Equal(t, nil, err)

	err = ListContribs(NewAppProject(filepath.Join(testEnv.currentDir, "myApp")), true, "")
	assert.Equal(t, nil, err)

}
