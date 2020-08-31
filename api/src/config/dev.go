// +build dev

package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func ConfigInit() {
	jsonFile, err := os.Open("./config/dev.json")
	log.Infof("Build: %s (Development)", BuildVersion)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Config)
}
