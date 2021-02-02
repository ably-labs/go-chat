package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ably/ably-go/ably"
	"github.com/ably/ably-go/ably/proto"
)

func main() {
	fmt.Println("Type your clientID")
	reader := bufio.NewReader(os.Stdin)
	clientId, _ := reader.ReadString('\n')
	clientId = strings.Replace(clientId, "\n", "", -1)

	opts := &ably.ClientOptions{
		AuthOptions: ably.AuthOptions{
			// If you have an Ably account, you can find
			// your API key at https://www.ably.io/accounts/any/apps/any/app_keys
			Key: "INSERT_YOUR_API_KEY_HERE",
		},
		ClientID: clientId,
		// NoEcho:   true, // Uncomment to stop messages you send from being sent back
	}

	fmt.Println("You can now send messages!")

	// Connect to Ably using the API key and ClientID specified above
	client, err := ably.NewRealtimeClient(opts)
	if err != nil {
		panic(err)
	}

	// Connect to the Ably Channel with name 'chat'
	channel := client.Channels.Get("chat")

	// Enter the Presence set of the channel
	channel.Presence.Enter("")

	getHistory(*channel)

	go subscribeToPresence(*channel)

	go subscribe(*channel)

	// Start the goroutine to allow for publishing messages
	publishing(*channel)
}

func subscribeToPresence(channel ably.RealtimeChannel) {
	sub, err := channel.Presence.Subscribe()
	if err != nil {
		panic(err)
	}

	for msg := range sub.PresenceChannel() {
		if msg.State == proto.PresenceEnter {
			fmt.Printf("%v has entered the chat\n", msg.ClientID)
		} else if msg.State == proto.PresenceLeave {
			fmt.Printf("%v has left the chat\n", msg.ClientID)
		}
	}
}

func getHistory(channel ably.RealtimeChannel) {
	// Before subscribing for messages, check the channel's
	// History for any missed messages. By default a channel
	// will keep 2 minutes of history available, but this can
	// be extended to 48 hours
	page, err := channel.History(nil)
	for ; err == nil; page, err = page.Next() {
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
