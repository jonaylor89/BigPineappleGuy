
package main

import (
	"log"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/dghubble/go-twitter/twitter"
)

const (
	credsFile = "creds.yml"
)

type creds struct {
	ConsumerKey string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
}

func getCreds() *creds {

	c := &creds{}

	// Read File
	yamlFile, err := ioutil.ReadFile(credsFile)
	if err != nil {
		log.Printf("[ERROR] yamlFile.Get err   #%v ", err)
	}

	// Unmarshall Creds
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	return c
}

func main() {

	creds := getCreds()

	if creds.ConsumerKey == "" ||  creds.ConsumerSecret == "" {
		log.Fatal("Application Access Token required")
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     creds.ConsumerKey,
		ClientSecret: creds.ConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// config := oauth1.NewConfig("consumerKey", "consumerSecret")
	// token := oauth1.NewToken("accessToken", "accessSecret")
	// httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	// jclient := twitter.NewClient(httpClient)
	
}