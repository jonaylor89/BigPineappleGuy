
package main

import (
	"log"
	"io/ioutil"
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
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	// Unmarshall Creds
	err = yaml.Unmarshal(yamlFile, v)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	return v.Victims
}