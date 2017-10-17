package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	// "time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
)

var AfkList []AfkUser = []AfkUser{}
var CustomCmdPrefix string = "-"

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

func isAdmin(user *discordgo.User) bool {
	// This is where i'm gonna add a buncha logic
	return true
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// ********* ADMIN COMMANDS **********

	if strings.HasPrefix(m.Content, CMD_PREFIX+"settag ") {
		// if isAdmin(m.Author) {
		// 	cmd_admin_addtag(s, m, *CustomCmdPrefix)
		// }
	}

	// ******** PROFILE COMMANDS *********

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
	if m.Content == CMD_PREFIX+"stats" || m.Content == CMD_PREFIX+"st" ||
		strings.HasPrefix(m.Content, CMD_PREFIX+"stats ") || strings.HasPrefix(m.Content, CMD_PREFIX+"st ") {
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

	// AFK command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"afk ") || m.Content == CMD_PREFIX+"afk" {
		log.Printf("AFKlist is currently this: %+v", AfkList)
		AfkList = cmd_afk_set(s, m, AfkList[0:])
		log.Printf("AFKlist now this: %v", AfkList)
	}
	AfkList = cmd_afk_check(s, m, AfkList)

	// Give cookie
	if strings.HasPrefix(m.Content, CMD_PREFIX+"cookie ") || strings.HasPrefix(m.Content, CMD_PREFIX+"c ") {
		cmd_cookie(s, m)
	}

	// ******** RP COMMANDS *********

	// Dice command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"roll ") || strings.HasPrefix(m.Content, CMD_PREFIX+"r ") {
		cmd_roll(s, m)
	}

	// Help item command
	if m.Content == CMD_PREFIX+"rpihelp" {
		rpcmd_item_help(s, m)
	}

	// Give item command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"giveitem") || strings.HasPrefix(m.Content, CMD_PREFIX+"gi") {
		rpcmd_item_give(s, m)
	}

	// Check inventory command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"rpinventory") || strings.HasPrefix(m.Content, CMD_PREFIX+"rpinv") {
		rpcmd_item_check(s, m)
	}

	// Delete own item command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"discarditem") || strings.HasPrefix(m.Content, CMD_PREFIX+"di") {
		rpcmd_item_discard(s, m)
	}

	// ******** FUNSIES COMMANDS *********

	if strings.Contains(m.Content, "well done") {
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank you!")
		if err != nil {
			fmt.Println("well done: Channel msg send unsuccessful")
			log.Printf("\n%v\n", err)
		}
	}

	// ******** DEV COMMANDS *********

	// Show information about the message sent
	if m.Content == CMD_PREFIX+"thism" {
		cmd_thisM(s, m)
	}

	// Changed as i dev
	if m.Content == CMD_PREFIX+"test" {
		m.Content = "$gi <@" + m.Author.ID + "> name:A dull, leather book - desc:thinger, alright? - weight: 0.1kg"
		rpcmd_item_give(s, m)
	}

	if m.Content == CMD_PREFIX+"test2" {
		m.Content = "$gi <@" + m.Author.ID + "> name:A pretty, preserved white flower - desc: A white clover plucked from the ground, and then dried for preservation - weight: 0.02kg"
		rpcmd_item_give(s, m)
	}

}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("A message has been deleted")
}
