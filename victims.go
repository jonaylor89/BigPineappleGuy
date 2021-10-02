
package main

import (
	"os"
	"log"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v2"
  )

const (
	victimsFile = "victims.yml"
)

type victims struct {
	Victims []string `yaml:"Victims"`
}

func getVictims() []string {

	v := &victims{}

	// Read File
	yamlFile, err := ioutil.ReadFile(victimsFile)
	if err == nil {

		// Unmarshall Creds
		err = yaml.Unmarshal(yamlFile, v)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}

	} else {

		// Fall back to using environment variables if no victims.yml file exists
		envVictims := os.Getenv("VICTIMS")
		if len(envVictims) > 0 {
			v.Victims = strings.Split(envVictims, ",")
		} else {
			// Fails as if no victims were provided (file or env)
			log.Fatalf("[ERROR] %v", err)
		}

	}



	return v.Victims
}
