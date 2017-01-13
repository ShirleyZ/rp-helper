package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"./emotes"
	"./profile"

	"github.com/ShirleyZ/godice"
	"github.com/bwmarrin/discordgo"
)

func cmd_credit(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing cmd: credit")
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

	err = profile.AddCredits(m.Author.ID, receiverName, amount)
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

func cmd_earnCredits(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing cmd: earn Credits")

	rand.Seed(time.Now().UTC().UnixNano())
	amount := rand.Intn(9) + 1
	log.Printf("\nEarning cash money: %v", amount)
	_ = profile.AddCredits(m.Author.ID, m.Author.Username, amount)
}

func cmd_register(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing cmd: register")
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf("register cannot find channel with that id")
		log.Printf("\n%v\n", err)
	}

	log.Println("Getting response channel...")
	sendToThis := m.ChannelID

	if channel.IsPrivate == false {
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("Unable to create private channel")
			log.Printf("\n%v\n", err)
		}
		sendToThis = channel.ID
	}
	log.Println("Done")
	log.Println("Calling profile.RegisterUser...")
	userInfo, err := profile.RegisterUser(m.Author.ID, m.Author.Username)
	log.Println("Done")

	if err != nil && err.Error() == profile.ERR_USEREXISTS {
		log.Println("Handling: User already exists")
		_, err = s.ChannelMessageSend(sendToThis, "You appear to have registered with us already. Please say $stats to check your details.")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
	} else if err != nil {
		log.Println("Handling: Other errors")
		// handle those errors yo
		_, err = s.ChannelMessageSend(sendToThis, "Oh dear! Something has gone wrong. Would you like to try and $register again?")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
		log.Println("Register: Something went wrong")
		log.Printf("\nRegistration Error:\n%v\n", err)
	} else {
		log.Println("Handling: New user created")
		log.Printf("\nUserInfo\n%+v\n", userInfo)
		// message := emotes.LookUpThis(content)
		_, err = s.ChannelMessageSend(sendToThis, "There we are! We've got your papers all set up now. Take a look")
		if err != nil {
			log.Printf("\nRegistration Error:\n%v\n", err)
		}
	}

	log.Println("=== End command")
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

func cmd_setProfile(s *discordgo.Session, m *discordgo.MessageCreate) {
	// get the stuff
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf("register cannot find channel with that id")
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

	if len(m.Content) <= len(CMD_PREFIX+"setprofile") {
		_, err = s.ChannelMessageSend(sendToThis, "Invalid format. Try $setprofile *<your text here>*")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	} else {
		profileBody := m.Content[len(CMD_PREFIX+"setprofile "):]
		result, err := profile.SetProfile(m.Author.ID, profileBody)
		if err != nil {
			log.Printf("\n%v\n", err)
		}
		log.Printf("\nResult\n%s", result)
		_, err = s.ChannelMessageSend(sendToThis, "There we are. I've updated your record with your new information.")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	}
}

func cmd_stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	var checkUser string
	if m.Content != CMD_PREFIX+"stats" {
		_, err := s.ChannelMessageSend(m.ChannelID, "This feature is currently disabled")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
		// checkUser = m.Content[len(CMD_PREFIX+"stats "):]
	} else if m.Content == CMD_PREFIX+"stats" {
		checkUser = m.Author.Username

		data, err := profile.CheckStats(m.Author.ID, checkUser)
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
			content := "```Markdown\n# == " + parsed.Title + " " + m.Author.Username + " == #\n"
			content += "* Credits: " + strconv.Itoa(parsed.Credits) + "\n"
			content += "* Profile: \n> " + parsed.Profile + "\n```"
			message := emotes.LookUpThis(content)
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Printf("\n%v\n", err)
			}
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
