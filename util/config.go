package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	prefFile = ".flogocli"
)

type Config struct {
}

func GetConfig() error {

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(currDir, prefFile)); err == nil {

		// local
	} else if _, err := os.Stat(filepath.Join(usr.HomeDir, prefFile)); err == nil {
		// home
	} else {
		//new
	}

	return nil
}

func loadConfig(configFile string) (map[string]interface{}, error) {

	jsonFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := make(map[string]interface{})
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
