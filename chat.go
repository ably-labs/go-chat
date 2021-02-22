package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ably/ably-go/ably"
	"github.com/ably/ably-go/ably/proto"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

  var username string

  // If no username specified, ask for one
  if len(os.Args) < 2 {
		fmt.Println("Type your username")
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.Replace(username, "\n", "", -1)
  } else {
  	username = os.Args[1]
  }

	opts := &ably.ClientOptions{
		AuthOptions: ably.AuthOptions{
			// If you have an Ably account, you can find
			// your API key at https://www.ably.io/accounts/any/apps/any/app_keys
			Key: os.Getenv("ABLY_KEY"),
		},
		ClientID: username,
		// NoEcho:   true, // Uncomment to stop messages you send from being sent back
	}

	// Connect to Ably using the API key and ClientID specified above
	client, err := ably.NewRealtimeClient(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println("You can now send messages!")

	// Connect to the Ably Channel with name 'chat'
	channel := client.Channels.Get("chat")

	// Enter the Presence set of the channel
	channel.Presence.Enter("")

	getHistory(*channel)

  // Subscribe to messages and presence messages
	go subscribe(*channel)
	go presenceSubscribe(*channel)

	// Start the goroutine to allow for publishing messages
	publishing(*channel)
}

func getHistory(channel ably.RealtimeChannel) {
	// Before subscribing for messages, check the channel's
	// History for any missed messages. By default a channel
	// will keep 2 minutes of history available, but this can
	// be extended to 48 hours
	page, err := channel.History(nil)
	for ; err == nil && page != nil; page, err = page.Next() {
		for _, msg := range page.Messages() {
			fmt.Printf("Previous message from %v: '%v'\n", msg.ClientID, msg.Data)
		}
	}
}

func subscribe(channel ably.RealtimeChannel) {
	// Initiate a subscription to the channel
	sub, err := channel.Subscribe()
	if err != nil {
		panic(err)
	}

	// For each message we receive from the subscription, print it out
	for msg := range sub.MessageChannel() {
		fmt.Printf("Received message from %v: '%v'\n", msg.ClientID, msg.Data)
	}
}

func presenceSubscribe(channel ably.RealtimeChannel) {
	// Subscribe to presence events (people entering and leaving) on the channel
	presenceSub, presenceErr := channel.Presence.Subscribe()
	if presenceErr != nil {
		panic(presenceErr)
	}

	for msg := range presenceSub.PresenceChannel() {
		if msg.State == proto.PresenceEnter {
			fmt.Printf("%v has entered the chat\n", msg.ClientID)
		} else if msg.State == proto.PresenceLeave {
			fmt.Printf("%v has left the chat\n", msg.ClientID)
		}
	}
}

func publishing(channel ably.RealtimeChannel) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		// Publish the message typed in to the Ably Channel
		res, err := channel.Publish("message", text)
		// await confirmation that message was received by Ably
		if err = res.Wait(); err != nil {
			panic(err)
		}
	}
}
