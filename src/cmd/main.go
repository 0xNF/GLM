package main

import (
	"fmt"
	"log"

	conf "github.com/0xNF/glm/src/internal/conf"
)

func main() {
	config := checkLoadConfig()
}

func checkLoadConfig() *conf.GLMConfig {

	var config *conf.GLMConfig = &conf.GLMConfig{}

	ok, err := conf.DoesConfigExist()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config file: %s", err))
	}
	if !ok {
		log.Println("local config file did not exist. Creating new one.")
		config, err = conf.MakeDefaultConfigFile()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to save new config file.: %s", err))
		}
		ok, err = conf.DoesConfigExist()
		if err != nil {
			log.Fatal(fmt.Sprintf("Created new config file, but it wasn't available to read: %s", err))
		}
	} else {
		config, err = conf.ReadConfigFile()
		if err != nil {
			log.Fatal(fmt.Sprintf("Found config file, but failed to load it: %s", err))
		}
	}

	return config
}
