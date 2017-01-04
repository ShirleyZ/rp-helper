package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	// "time"

	"./emotes"
	"./profile"

	"github.com/ShirleyZ/godice"
	"github.com/bwmarrin/discordgo"
)

func cmd_credit(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse arguments [who, amount]
	args := strings.Split(m.Content, " ")
	amount, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		log.Println("Error: Unable to parse amount")
		log.Printf("\n%v\n", err)
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

func cmd_register(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Executing cmd: register")
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf(CMD_PREFIX + "register cannot find channel with that id")
		log.Printf("\n%v\n", err)
	}

	sendToThis := m.ChannelID

	if channel.IsPrivate == false {
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("Unable to create private channel")
			log.Printf("\n%v\n", err)
		}
		sendToThis = channel.ID
	}
	userInfo, err := profile.RegisterUser(m.Author.Username)

	if err != nil && err.Error() == profile.ERR_USEREXISTS {
		_, err = s.ChannelMessageSend(sendToThis, "You appear to have registered with us already. Please say $stats to check your details.")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
		log.Println("Register: User already exists")
	} else if err != nil {
		// handle those errors yo
		_, err = s.ChannelMessageSend(sendToThis, "Oh dear! Something has gone wrong. Would you like to try and $register again?")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
		log.Println("Register: Something went wrong")
		log.Printf("\nRegistration Error:\n%v\n", err)
	} else {
		log.Printf("\nUserInfo\n%+v\n", userInfo)
		// message := emotes.LookUpThis(content)
		_, err = s.ChannelMessageSend(sendToThis, "There we are! We've got your papers all set up now. Take a look")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
	}
}

func cmd_roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Print("Rolling")
	result, err := dice.Roll(m.Content[len(CMD_PREFIX+"roll"):])
	result = m.Author.Username + " " + result
	if err == nil {
		_, err = s.ChannelMessageSend(m.ChannelID, result)
		if err != nil {
			fmt.Println("!!roll: Channel msg send unsuccessful")
			log.Printf("\n%v\n", err)
		}
	}
}

func cmd_stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	var checkUser string
	if m.Content == CMD_PREFIX+"stats" {
		checkUser = m.Author.Username
	} else {
		checkUser = m.Content[len(CMD_PREFIX+"stats "):]
	}

	data, err := profile.CheckStats(checkUser)
	if err != nil && err.Error() == "No user with that name found" {
		_, err = s.ChannelMessageSend(m.ChannelID, "I don't seem to have your record on file. Please message me to $register.")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	} else if err != nil {
		log.Println("Error checking user stats")
		log.Printf("\n%v\n", err)
	} else {
		parsed := UserProfile{}
		err = json.Unmarshal([]byte(data), &parsed)
		if err != nil {
			log.Printf("\n%v\n", err)
		}
		content := "```Markdown\n# == " + parsed.Title + " " + parsed.Username + " == #\n"
		content += "* Credits: " + strconv.Itoa(parsed.Credits) + "\n"
		content += "* Profile: \n> " + parsed.Profile + "\n```"
		message := emotes.LookUpThis(content)
		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	}
}

func cmd_thisM(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Executing cmd: !thism")
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf(CMD_PREFIX + "thism cannot find channel with that id")
		log.Printf("\n%v\n", err)
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
			log.Printf("\n%v\n", err)
		}
		msg += "Server: (" + channel.GuildID + ")\n"
		// msg += fmt.Sprintf("Author roles: %+v \n", member.Roles)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		fmt.Println("!thism: Channel msg send unsuccessful")
		log.Printf("\n%v\n", err)
	}
}
