package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func IsRemote(path string) bool {
	return strings.HasPrefix(path, "http")
}

func LoadRemoteFile(sourceURL string) (string, error) {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func LoadLocalFile(path string) (string, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CopyFile(srcFile, destFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
