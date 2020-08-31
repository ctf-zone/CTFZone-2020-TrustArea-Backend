// +build !dev

package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func ConfigInit() {
	jsonFile, err := os.Open("./config/prod.json")
	log.Infof("Build: %s (Production)", BuildVersion)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Config)
}
