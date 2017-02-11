package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type UserProfile struct {
	Id       string
	Username string
	Credits  int
	Profile  string
	Title    string
	Cookies  int
}

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

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// ******** DEV COMMANDS *********

	// Show information about the message sent
	if m.Content == CMD_PREFIX+"thism" {
		cmd_thisM(s, m)
	}

	// If the message is "!test" send the message to server-chatter
	if m.Content == CMD_PREFIX+"test" {
		// dice.Roll("1d5")
	}

	// ******** RP COMMANDS *********

	cmd_earnCredits(s, m)

	// Register a new account
	if m.Content == CMD_PREFIX+"register" {
		cmd_register(s, m)
	}

	// Tags
	if strings.HasPrefix(m.Content, CMD_PREFIX+"t ") {
		cmd_tags(s, m)
	}

	// Halp
	if strings.HasPrefix(m.Content, CMD_PREFIX+"help") {
		cmd_help(s, m)
	}

	// Stats
	// if strings.HasPrefix(m.Content, CMD_PREFIX+"stats") || strings.HasPrefix(m.Content, CMD_PREFIX+"st") {
	if m.Content == CMD_PREFIX+"stats" || m.Content == CMD_PREFIX+"st" ||
		strings.HasPrefix(m.Content, CMD_PREFIX+"stats") || strings.HasPrefix(m.Content, CMD_PREFIX+"st") {
		cmd_stats(s, m)
	}

	// Add credits command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"credit") {
		if m.Author.Username == "Dintay" {
			cmd_credit(s, m)
		}
	}

	// Set profile
	if strings.HasPrefix(m.Content, CMD_PREFIX+"profile ") || strings.HasPrefix(m.Content, CMD_PREFIX+"p ") {
		cmd_setProfile(s, m)
	}

	// Dice command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"roll ") || strings.HasPrefix(m.Content, CMD_PREFIX+"r ") {
		cmd_roll(s, m)
	}

	// Give cookie
	if strings.HasPrefix(m.Content, CMD_PREFIX+"cookie ") || strings.HasPrefix(m.Content, CMD_PREFIX+"c ") {
		cmd_cookie(s, m)
	}

	if strings.Contains(m.Content, "well done") {
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank you!")
		if err != nil {
			fmt.Println("well done: Channel msg send unsuccessful")
			log.Printf("\n%v\n", err)
		}
	}

}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("A message has been deleted")
}
