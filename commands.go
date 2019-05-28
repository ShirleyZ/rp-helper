package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"./emotes"
	"./profile"

	"github.com/ShirleyZ/godice"
	"github.com/bwmarrin/discordgo"
)

const MAX_COOKIES_GIVEN = 9001

type AfkUser struct {
	Id       string
	Username string
	Reason   string
	AfkStart time.Time
}

/* Function takes a string and takes every word
 * separated by a space to be a new argument
 */
func get_longParam(cmd string) []string {
	var params []string
	for i := range cmd {
		if cmd[i] == ' ' {
			params[0] = cmd[:i]
			params[1] = cmd[i+1:]
			return params
		}
	}
	return params
}

/* addTag command requirements
 * - Tag has to be one word
 */
func cmd_admin_addtag(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get the params
	// params := get_longParam(m.Content)
}

func cmd_admin_setCustomPrefix(s *discordgo.Session, m *discordgo.MessageCreate, customPrefix *string) {

}

func cmd_afk_check(s *discordgo.Session, m *discordgo.MessageCreate, afkList []AfkUser) []AfkUser {
	log.Println("=== Executing cmd: afk checking")

	if !strings.HasPrefix(m.Content, CMD_PREFIX) {

		// ********** ALERTING ABT AFK USERS *************

		// Check if any are afk
		dudesAfk := []int{}
		authorActive := false
		authorIndex := 0

		for afkIndex, afkVal := range afkList {
			for _, user := range m.Mentions {
				// log.Printf("AFK %s %s compared against MENTIONED %s", afkVal.Id, afkVal.Username, user.ID)
				if afkVal.Id == user.ID {
					dudesAfk = append(dudesAfk, afkIndex)
				}
			}
			if afkVal.Id == m.Author.ID {
				authorActive = true
				authorIndex = afkIndex
			}
		}

		// log.Printf("These pinged people are afk: %+v", dudesAfk)

		// If any show up, create message for afk users
		// save ids instead, so you can access all the info
		for _, afIn := range dudesAfk {
			currUser := afkList[afIn]

			message := "<@" + currUser.Id + "> is AFK. "
			if currUser.Reason != "No reason given" {
				message += "Reason: "
			}
			message += currUser.Reason

			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Printf("\n%v\n", err)
			}
		}

		// ********** REMOVING AFK STATUS *************
		if authorActive {
			afkList = append(afkList[:authorIndex], afkList[authorIndex+1:]...)
			_, err := s.ChannelMessageSend(m.ChannelID, "Welcome back <@"+m.Author.ID+">. AFK status removed.")
			if err != nil {
				log.Printf("\n%v\n", err)
			}
		}
	}

	return afkList
}

func cmd_afk_set(s *discordgo.Session, m *discordgo.MessageCreate, afkList []AfkUser) []AfkUser {
	fmt.Println("=== Executing cmd: afk setting")

	afkReason := "No reason given"
	if m.Content == CMD_PREFIX+"afk" {
		// No reason, do nothing
	} else {

		// Parse it to remove any USER PINGS -SIDEEYES ZAZ-
		betterMessage := m.ContentWithMentionsReplaced()

		afkReason = betterMessage[len(CMD_PREFIX+"afk "):]

	}
	newAfk := AfkUser{m.Author.ID, m.Author.Username, afkReason, time.Now()}
	// log.Printf("%+v", newAfk)
	log.Printf("%s is afk bec %s", newAfk.Username, newAfk.Reason)

	// Look through current list to see if you need to override old afk
	found := false
	index := 0
	for i, elem := range afkList {
		if elem.Id == m.Author.ID {
			found = true
			index = i
			elem.Reason = afkReason
		}
	}

	if found == false {
		afkList = append(afkList, newAfk)
	} else {
		afkList = append(afkList[:index], afkList[index+1:]...)
		afkList = append(afkList, newAfk)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "AFK set for <@"+newAfk.Id+"> Reason: "+newAfk.Reason)
	if err != nil {
		log.Printf("\n%v\n", err)
	}
	return afkList
}

func cmd_cookie(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing cmd: cookie")
	// Parse command for target
	r, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Println("Regexp unsuccessful init")
	}
	recipientId := r.Find([]byte(m.Content))

	giverId := m.Author.ID
	amount := 1
	args := strings.Split(m.Content, " ")

	// Invalid command params
	if len(args) != 3 {
		log.Println("Incorrect parameters")
	} else {
		amount, _ = strconv.Atoi(args[2])

		// Invalid input params
		if string(recipientId) == "" {
			log.Println("No Recipient")
		} else if amount <= 0 {
			log.Println("Less than 0 amount")
		} else if amount > MAX_COOKIES_GIVEN {
			log.Println("More than max amount")
		} else {
			log.Printf("Sender: %s Receiver: %s Amount: %s", giverId, recipientId, amount)

			err = profile.GiveCookie(string(giverId), string(recipientId), amount)
			if err != nil {
				_, err = s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, "Success")
			}
		}
	}

}

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

	log.Printf("Giving %+v %d credits", receiverName, amount)

	err = profile.AddCredits(m.Author.ID, receiverName, amount)
	if err != nil {
		log.Println("Unsuccessful attemptt o add")
		// Do nothing
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
	log.Printf("%v earned %v credits", m.Author.Username, amount)
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

	if channel.Type == discordgo.ChannelTypeDM {
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

func cmd_help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == CMD_PREFIX+"help" {
		message := msg_help()
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println("!!roll: Channel msg send unsuccessful")
			log.Printf("\n%v\n", err)
		}
	} else {
		// Parse for what you need help w/
	}
}

func cmd_roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Print("Rolling")

	cmdAlias := ""
	if strings.HasPrefix(m.Content, CMD_PREFIX+"roll ") {
		cmdAlias = "$roll"
	} else if strings.HasPrefix(m.Content, CMD_PREFIX+"r ") {
		cmdAlias = "$r"
	}

	result, err := dice.Roll(m.Content[len(cmdAlias):])
	if err != nil {
		// Do nothing
		fmt.Println("!!roll: Error executing dice.Roll")
		log.Printf("\n%v\n", err)
	} else {
		result = m.Author.Username + " " + result
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

	cmdAlias := ""
	if strings.HasPrefix(m.Content, CMD_PREFIX+"profile ") {
		cmdAlias = "$profile"
	} else if strings.HasPrefix(m.Content, CMD_PREFIX+"p ") {
		cmdAlias = "$p"
	}

	if channel.Type != discordgo.ChannelTypeDM {
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("Unable to create private channel")
			log.Printf("\n%v\n", err)
		}
		sendToThis = channel.ID
	}

	if len(m.Content) <= len(cmdAlias) {
		// Do nothing
	} else {
		profileBody := m.Content[len(cmdAlias+" "):]
		result, err := profile.SetProfile(m.Author.ID, profileBody)
		if err != nil {
			log.Printf("\n%v\n", err)
		}
		log.Printf("\nResult\n%s", result)
		_, err = s.ChannelMessageSend(sendToThis, "There we are. I've updated your record with your new information.")

		userCard, err := msg_profile(m, nil)
		if err != nil {
			log.Printf("\n%v\n", err)
		}
		_, err = s.ChannelMessageSend(sendToThis, userCard)
	}
}

func cmd_stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error
	var content string
	wrongParams := false

	args := strings.Split(m.Content, " ")

	if len(args) > 2 {
		// Do nothing
		// Invalid params
	} else {
		var userId []byte = nil
		// If you're looking up for someone else
		if len(args) == 2 {
			r, err := regexp.Compile("[0-9]+")
			if err != nil {
				log.Println("Regexp unsuccessful init")
			}
			userId = r.Find([]byte(m.Content))
			if userId == nil {
				wrongParams = true
			}
		}

		if wrongParams == true {
			_, err = s.ChannelMessageSend(m.ChannelID, "User incorrectly specified")
			if err != nil {
				log.Printf("\n%v\n", err)
			}
		} else {
			content, err = msg_profile(m, userId)

			if err != nil {
				log.Println("There's an error bud")
				log.Printf("\n%v\n", err)
				_, err = s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					log.Printf("\n%v\n", err)
				}
			} else {
				message := emotes.LookUpThis(content)
				_, err = s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Printf("\n%v\n", err)
				}

			}
		}

	}
}

func cmd_tags(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == CMD_PREFIX+"t bodyhorror" {
		_, err := s.ChannelMessageSend(m.ChannelID, "**BODY HORROR WARNING**")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	} else if m.Content == CMD_PREFIX+"t confirmed" {
		_, err := s.ChannelMessageSend(m.ChannelID, "**CONFIRMED**")
		if err != nil {
			log.Printf("\n%v\n", err)
		}
	} else if m.Content == CMD_PREFIX+"t banned" {
		_, err := s.ChannelMessageSend(m.ChannelID, "**BANNED**")
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

	msg := ""
	// msg += "Sent: " + m.Timestamp + "\n"
	msg += "Channel: " + channel.Name + "(" + m.ChannelID + ")\n"
	msg += "Type: " + fmt.Sprintf("%i", channel.Type) + "\n"
	// if channel.IsPrivate == true {
	// 	msg += "**== PM details ==**\n"
	// 	msg += "Recipient: " + channel.Recipient.Username + "\n"
	// 	msg += "Author: " + m.Author.Username + "(" + m.Author.ID + ")\n"
	// 	log.Printf("\n%+v\n", m.Author)
	// } else {
	// 	msg += "**== Channel details ==**\n"
	// 	permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	// 	log.Printf("\n%v\n", permissions)
	// 	// member, err := s.GuildMember(channel.GuildID, m.Author.ID)
	// 	if err != nil {
	// 		fmt.Printf(CMD_PREFIX + "thism cannot find guild member")
	// 		log.Printf("\n%v\n", err)
	// 	}
	// 	msg += "Server: (" + channel.GuildID + ")\n"
	// 	// msg += fmt.Sprintf("Author roles: %+v \n", member.Roles)
	// }

	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		fmt.Println("!thism: Channel msg send unsuccessful")
		log.Printf("\n%v\n", err)
	}
}

func msg_help() string {
	message := "*Hello! Here are the list of things I am able to help you with!*"
	message += "\n```Markdown"
	message += "\n# == Commands List == #"
	message += "\n== User account"
	message += "\no $register       - create an account with Scrivener Nibb"
	message += "\no $stats @<user>  - check your stats"
	message += "\n                  - alias: $st"
	message += "\no $profile <text> - set your profile text"
	message += "\n                  - alias: $p"
	message += "\n\n== Funsies"
	message += "\no $cookie @<user> <amount> - Buy a cookie for the pinged user"
	message += "\n                  - Cookies cost 20ea"
	message += "\n                  - alias: $c"
	message += "\no $roll #d# <action> - roll to make an action eg roll 1d20 to party"
	message += "\n                  - Max 100 dice at a time"
	message += "\n                  - alias: $r"
	message += "\n\n== Features"
	message += "\no $rpihelp				- Gives a list of freeform inventory commands"
	message += "```"

	return message
}
func msg_profile(m *discordgo.MessageCreate, userId []byte) (string, error) {
	var err error
	var data string
	if userId == nil {
		data, err = profile.CheckStats(m.Author.ID)
	} else {
		data, err = profile.CheckStats(string(userId))
	}
	if err != nil {
		log.Printf("Errorrr: %s", err.Error())
		return "", err
	}
	userInfo := profile.UserData{}
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		log.Printf("%v", err)
		return "", nil
	}

	content := "```Markdown\n# == " + userInfo.Title + " " + strings.Title(userInfo.Username) + " == #\n"
	content += "* Credits: " + strconv.Itoa(userInfo.Credits) + "\n"
	content += "* Cookies: " + strconv.Itoa(userInfo.Cookies) + "\n"
	content += "# ==== #\n"
	content += "* Profile: \n> " + userInfo.Profile + "\n"
	content += "```"
	return content, nil
}

// func util_getGuildId(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
// 	channelInfo, err := s.Channel(m.ChannelID)
// 	if err != nil {
// 		log.Println("Cannot get server info")
// 		log.Printf("Error: %+v", err.Error())
// 		return "", err
// 	}
// 	guildId := channelInfo.GuildID
// 	return guildId, nil
// }
