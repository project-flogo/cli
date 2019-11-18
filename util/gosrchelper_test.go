package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindGoModPackageSrc(t *testing.T) {
	str, ver, err := FindGoModPackageSrc("github.com/project-flogo/core", "", true)
	if err != nil {
		fmt.Println("err:", err)
		t.FailNow()
	}

	fmt.Println("path: ", str)
	fmt.Println("ver: ", ver)
}


func TestFindOldPackageSrc(t *testing.T) {
	str, ver, err := FindOldPackageSrc("github.com/project-flogo/cli")
	if err != nil {
		fmt.Println("err:", err)
		t.FailNow()
	}

	fmt.Println("path: ", str)
	fmt.Println("ver: ", ver)
}


func TestFindGoModPackageSrcNotFound(t *testing.T) {
	_, _, err := FindGoModPackageSrc("github.com/project-blah/core", "", true)
	assert.True(t, IsPkgNotFoundError(err))
	fmt.Println("err: ", err)
}


func TestFindOldPackageSrcNotFound(t *testing.T) {
	_, _, err := FindOldPackageSrc("github.com/project-blah/core")
	assert.True(t, IsPkgNotFoundError(err))
	fmt.Println("err: ", err)
}
