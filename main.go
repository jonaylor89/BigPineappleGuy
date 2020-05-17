
package main

import (
	"os"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"math/rand"

	"github.com/dghubble/oauth1"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/clientcredentials"
	"github.com/dghubble/go-twitter/twitter"
)

var victims = []string{
	"A_Chris_Kahuna",
}

func main() {

	creds := getCreds()

	if creds.ConsumerKey == "" ||  creds.ConsumerSecret == "" || creds.AccessToken == "" || creds.AccessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {

		choice := facts[rand.Intn(len(facts))]
		fmt.Println("[INFO] Tweet: ", tweet.Text)
		fmt.Println("[INFO] Pineapple Fact: ", choice)

		// Send a Tweet
		// tweet, resp, err := client.Statuses.Update(
		// 	choice, 
		// 	&StatusUpdateParams{ 
		// 		InReplyToStatusID: tweet.ID
		// 	},
		// )
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println("[INFO] DM: ", dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("[INFO] Event: %#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Follow:         victims,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}