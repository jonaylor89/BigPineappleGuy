package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
)

func main() {

	creds := getCreds()
	victims := getVictims()

	if creds.ConsumerKey == "" || creds.ConsumerSecret == "" || creds.AccessToken == "" || creds.AccessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	rand.Seed(time.Now().UnixNano())

	var ids = []string{}

	for _, name := range victims {

		// User Show
		user, _, err := client.Users.Show(&twitter.UserShowParams{
    		ScreenName: name,
		})
		if err != nil {
			fmt.Println("[ERROR] on ", name)
			continue
		}
	
		ids = append(ids, strconv.FormatInt(user.ID, 10))
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {

		// Ignore RTs
		if tweet.Retweeted {
			return
		}

		// Ignore Replies
		if tweet.InReplyToStatusID != 0 || tweet.InReplyToScreenName != "" || tweet.InReplyToUserIDStr != "" {
			return
		}

		choice := facts[rand.Intn(len(facts))]
		fmt.Println("[INFO] Tweet: ", tweet.Text)

		botResponse := fmt.Sprintf("@%s Pineapple Fact: %s", tweet.User.ScreenName, choice)

		// Reply to  Tweet
		reply, _, err := client.Statuses.Update(
			botResponse,
			&twitter.StatusUpdateParams{
				InReplyToStatusID: tweet.ID,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("[INFO] ", reply)
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
		Follow:        ids,
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
