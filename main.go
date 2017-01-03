package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ShirleyZ/godice"
	"github.com/bwmarrin/discordgo"

	"./emotes"
	"./profile"
)

type UserProfile struct {
	Username string
	Credits  int
	Profile  string
	Title    string
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
	if m.Content == CMD_PREFIX+"thism" {
		fmt.Println("Executing cmd: !thism")
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Printf(CMD_PREFIX + "thism cannot find channel with that id")
			log.Fatal(err)
		}

		msg := "Sent: " + m.Timestamp + "\n"
		msg += "Channel: " + channel.Name + "(" + m.ChannelID + ")\n"
		msg += "isPrivate: " + fmt.Sprintf("%t", channel.IsPrivate) + "\n"
		if channel.IsPrivate == true {
			msg += "**== PM details ==**\n"
			msg += "Recipient: " + channel.Recipient.Username + "\n"
			msg += "Author: " + m.Author.Username + "(" + m.Author.ID + ")\n"
			log.Printf("\n%+v\n", m.Author)
		} else {
			msg += "**== Channel details ==**\n"
			permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
			log.Printf("\n%v\n", permissions)
			// member, err := s.GuildMember(channel.GuildID, m.Author.ID)
			if err != nil {
				fmt.Printf(CMD_PREFIX + "thism cannot find guild member")
				log.Fatal(err)
			}
			msg += "Server: (" + channel.GuildID + ")\n"
			// msg += fmt.Sprintf("Author roles: %+v \n", member.Roles)
		}

		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("!thism: Channel msg send unsuccessful")
			log.Fatal(err)
		}
	}

	// If the message is "!test" send the message to server-chatter
	if m.Content == CMD_PREFIX+"test" {
		// dice.Roll("1d5")
	}

	// Register a new account
	if m.Content == CMD_PREFIX+"register" {
		fmt.Println("Executing cmd: register")
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Printf(CMD_PREFIX + "register cannot find channel with that id")
			log.Fatal(err)
		}

		if channel.IsPrivate == false {
			_, err = s.ChannelMessageSend(m.ChannelID, "Please message me privately, and we can sort out the necessary paperwork.")
			if err != nil {
				fmt.Println("register: Channel msg send unsuccessful")
				log.Fatal(err)
			}
		} else {
			userInfo, err := profile.RegisterUser(m.Author.Username)
			if err != nil {
				// handle those errors yo
				log.Println("Register: Something went wrong")
				log.Fatal(err)
			} else {

			}
		}
	}

	// see profile stats
	if strings.HasPrefix(m.Content, CMD_PREFIX+"stats") {
		var checkUser string
		if m.Content == CMD_PREFIX+"stats" {
			checkUser = m.Author.Username
		} else {
			checkUser = m.Content[len(CMD_PREFIX+"stats "):]
		}

		data, err := profile.CheckStats(checkUser)
		if err == "No user with that name found" {
			_, err = s.ChannelMessageSend(m.ChannelID, "I don't seem to have your record on file. Please message me to $register.")
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Println("Error checking user stats")
			log.Fatal(err)
		} else {
			parsed := UserProfile{}
			err = json.Unmarshal([]byte(data), &parsed)
			if err != nil {
				log.Fatal(err)
			}
			content := "```Markdown\n# == " + parsed.Title + " " + parsed.Username + " == #\n"
			content += "* Credits: " + strconv.Itoa(parsed.Credits) + "\n"
			content += "* Profile: \n> " + parsed.Profile + "\n```"
			message := emotes.LookUpThis(content)
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	// Add credits command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"credit") {
		// Parse arguments [who, amount]
		args := strings.Split(m.Content, " ")
		amount, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			log.Println("Error: Unable to parse amount")
			log.Fatal(err)
		}
		receiverName := ""
		if len(args) > 3 {
			for i := 1; i < len(args)-1; i++ {
				receiverName += args[i]
				if i != len(args)-1 {
					receiverName += " "
				}
			}
		} else {
			receiverName = args[1]
		}

		log.Printf("\nReceiver Name\n%+v\n", receiverName)
		log.Printf("\nAmount\n%d\n", amount)

		err = profile.AddCredits(receiverName, amount)
		if err != nil {
			log.Println("Unsuccessful attemptt o add")
			_, err = s.ChannelMessageSend(m.ChannelID, "UnSuccess")
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Success")
			if err != nil {
				fmt.Println("unable to send message")
			}

		}

	}

	// Dice command
	if strings.HasPrefix(m.Content, CMD_PREFIX+"roll") {
		fmt.Print("Rolling")
		result, err := dice.Roll(m.Content[len(CMD_PREFIX+"roll"):])
		result = m.Author.Username + " " + result
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
