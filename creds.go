
package main

import (
  "fmt"
	"os"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	credsFile = "creds.yml"
)

type creds struct {
	ConsumerKey string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
	AccessToken string `yaml:"AccessToken"`
	AccessSecret string `yaml:"AccessSecret"`
}

func getCreds() *creds {

	c := &creds{}

	// Read File
	yamlFile, err := ioutil.ReadFile(credsFile)
	if err != nil {
    fmt.Println("[INFO] No yml config, pulling from environment")
		c = &creds{
			ConsumerKey: os.Getenv("CONSUMER_KEY"),
			ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
			AccessToken: os.Getenv("ACCESS_TOKEN"),
			AccessSecret: os.Getenv("ACCESS_SECRET"),
		}
	} else {
		// Unmarshall Creds
		err = yaml.Unmarshal(yamlFile, c)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
	}

	return c
}
