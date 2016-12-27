package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ShirleyZ/godice"
	"github.com/bwmarrin/discordgo"
)

// const CHANNEL_SERVER_CHATTER = "259926279820804097"

// Variables used for command line parameters
var (
	Token string
	BotID string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}
	fmt.Println("=== Received my user details. They are: ")
	fmt.Println(u)

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageDelete)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	// fmt.Print(m)

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		fmt.Print("Executing cmd: ping")
		_, err := s.ChannelMessageSend(m.ChannelID, "Pongg!")
		if err != nil {
			fmt.Println("ping: channelmsgsend error")
			log.Fatal(err)
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		fmt.Println("Executing cmd: pong")
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			fmt.Print("pong: channelmsgsend error")
			log.Fatal(err)
		}
	}

	// Show information about the message sent
	if m.Content == "!thism" {
		fmt.Println("Executing cmd: !thism")
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Printf("!thism cannot find channel with that id")
			log.Fatal(err)
		}
		msg := "Sent: " + m.Timestamp + "\n"
		msg += "Channel: " + channel.Name + "(" + m.ChannelID + ")\n"
		msg += "isPrivate: " + fmt.Sprintf("%t", channel.IsPrivate) + "\n"
		if channel.IsPrivate == true {
			msg += "Recipient: " + channel.Recipient.Username
		}

		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("!thism: Channel msg send unsuccessful")
			log.Fatal(err)
			// return
		}
	}

	// If the message is "!test" send the message to server-chatter
	if m.Content == "!test" {
		// dice.Roll("1d5")
	}

	// Dice command
	fmt.Printf("\n%s\n", m.Content)
	if strings.HasPrefix(m.Content, "!!roll") {
		fmt.Print("Rolling")
		result, err := dice.Roll(m.Content)
		if err == nil {
			_, err = s.ChannelMessageSend(m.ChannelID, result)
			if err != nil {
				fmt.Println("!!roll: Channel msg send unsuccessful")
				log.Fatal(err)
			}
		}
	}

	if strings.Contains(m.Content, "well done") {
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank you!")
		if err != nil {
			fmt.Println("well done: Channel msg send unsuccessful")
			log.Fatal(err)
		}
	}

}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("A message has been deleted")
}
