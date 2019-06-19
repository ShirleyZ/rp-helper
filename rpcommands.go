package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "reflect"
	// "math/rand"
	"regexp"
	// "strconv"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	// "time"

	"github.com/bwmarrin/discordgo"
)

const API_ENDPOINT = "http://localhost:8080/api/"
const ITEM_PROP_LIMIT = 15

type UserItemData struct {
	UserId string              `bson:"userid" json:"userid"`
	Items  []map[string]string `bson:"items" json:"items"`
}

// ********** Bot Administration System ********** //
// - Feature toggles
// - Settings for commands/systems

// ********** Freeform Inventory System ********** //
// - Admin is able to remove someone's item
// - User is able to transfer their own item to someone else
// - User is able to remove their own items
// - User can see their own inventory
// --- Nice to haves
// - Items are to be displayed on user's profile
// - Items to keep history of ownership
// - Users can see other user's inventory
// - Users can hide items in their inventory (in a bag?)
// - Admin is able to remove someone's item with a reason
// - Admin logging functionality
// --- A/C
// - Items need to have a name, description, date created field
func rpcmd_item_check(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: checkInventory")
	userId := m.Author.ID

	// Getting guild id
	guildId, err := util_getGuildId(s, m)
	if err != nil {
		log.Println("Cannot get channelid")
		log.Printf("Error: %+v", err.Error())
		return
	}

	// Do a lookup of their own stuff
	endpoint := API_ENDPOINT + "rpcmd/item/check/"
	sendBody, err := url.ParseQuery("userid=" + userId + "&guildid=" + guildId)
	if err != nil {
		log.Println("Cannot parse query")
		log.Printf("Error: %+v", err.Error())
		return
	}

	// Call endpoint to check inv
	fmt.Printf("\nHitting endpoint: %s\n", endpoint)
	resp, err := http.PostForm(endpoint, sendBody)
	if err != nil {
		log.Println("didn't success the post")
		log.Printf("%+v", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("didn't suceess reading the body")
		log.Fatal(err)
	}

	// log.Printf("Stringed Body: %+v", string(body))

	if body == nil {
		return
	}

	userInventory := UserItemData{}
	err = json.Unmarshal(body, &userInventory)
	if err != nil {
		log.Println("Error: Couldn't unmarshal")
		log.Printf("%+v", err.Error())
		return
	}

	message := "```Markdown\n"
	message += "*== " + m.Author.Username + "'s Inventory ==*\n\n"

	if len(userInventory.Items) == 0 {
		message += "There is nothing here"
	} else {
		for key, value := range userInventory.Items {
			currItem := value
			log.Printf("%s: %+v ", key, currItem)
			message += "# [" + currItem["itemid"] + "] " + currItem["name"] + " #\n"
			for propName, propValue := range currItem {
				if (propName != "itemid") && (propName != "name") {
					message += "- " + propName + ": " + propValue + "\n"
				}
			}
		}
	}

	message += "```"
	// Different functions for
	// brief: id, name only
	// verbose: all fields
	// overview: id, name, certain fields **Nice to have

	_, err = s.ChannelMessageSend(m.ChannelID, message)
}

func rpcmd_item_discard(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: discardItem")
	// User removes their own item

	numParams := strings.Count(m.Content, " ")
	if numParams < 1 {
		log.Println("Error: Not enough params")
		return
	}
	// Get user id
	userId := m.Author.ID
	guildId, err := util_getGuildId(s, m)
	if err != nil {
		log.Println("Error: " + err.Error())
		return
	}

	// Get item name or id
	cmdParam := ""
	if strings.HasPrefix(m.Content, "$rpid ") {
		cmdParam = m.Content[len("$rpid "):]
	} else if strings.HasPrefix(m.Content, "$rpidiscard ") {
		cmdParam = m.Content[len("$rpidiscard "):]
	}
	cmdParam = strings.ToUpper(cmdParam)
	log.Printf("cmdParam: %s", cmdParam)

	itemId := strings.Trim(cmdParam, " ")
	log.Printf("itemId: %s", itemId)
	// Hit endpoint
	sendBody, err := url.ParseQuery("userid=" + userId + "&itemid=" + itemId + "&guildid=" + guildId)
	if err != nil {
		log.Println("Cannot parse query")
		log.Printf("Error: %+v", err.Error())
		return
	}
	url := API_ENDPOINT + "rpcmd/item/discard/"
	fmt.Printf("\nHitting endpoint: %s\n", url)
	resp, err := http.PostForm(url, sendBody)
	if err != nil {
		log.Println("didn't success the post")
		log.Printf("%+v", err)
		return
	}
	defer resp.Body.Close()

	message := ""
	if resp.StatusCode == 200 { // OK
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("resp string: %s", bodyString)
		if bodyString == "OK" {
			message = "Item " + itemId + " has been discarded"
		} else if bodyString != "" {
			message = bodyString
		}
	}
	log.Printf("resp body: %+v", resp.Body)

	if message != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, message)
	}
}

func rpcmd_item_give(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: giveItem")

	// Check there are enough params
	numParams := strings.Count(m.Content, " ")
	if numParams < 2 {
		log.Println("Error: Not enough params")
		return
	}

	params, err := util_getParams_itemGive(m.Content, true)
	noMentions := m.ContentWithMentionsReplaced()
	log.Printf("NO MNTIONS: %s", noMentions)
	log.Printf("Params: %+v", params)

	recipientUsername := ""
	if len(m.Mentions) > 0 {
		recipientUsername = m.Mentions[0].Username
	}

	if err != nil {
		log.Println("Error: problem getting params for item give")
		return
		// TODO: handle error
	}
	fmt.Printf("\nparams:\n%+v\n", params)
	// TODO: bot-end checking of param

	guildId, err := util_getGuildId(s, m)
	if err != nil {
		log.Println("Cannot get channelid")
		log.Printf("Error: %+v", err.Error())
		return
	}

	sendBody, err := url.ParseQuery("userid=" + params["user"] + "&itemparams=" + params["item"] + "&guildid=" + guildId)
	if err != nil {
		log.Println("Cannot parse query")
		log.Printf("Error: %+v", err.Error())
		return
	}
	url := API_ENDPOINT + "rpcmd/item/give/"
	fmt.Printf("\nHitting endpoint: %s\n", url)
	thing, err := http.PostForm(url, sendBody)
	if err != nil {
		log.Println("didn't success the post")
		log.Printf("%+v", err)
		return
	}

	log.Printf("thing: %+v", thing)
	_, err = s.ChannelMessageSend(m.ChannelID, "Item given to "+recipientUsername)
}

func rpcmd_item_help(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: helpItem")
	message := msg_rpcmd_help()
	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}

// ********** Statistics System ********** //
// - Initialises default set of stats when feature is turned on
// - Admin is able to set server wide custom stats
// - Admin is able to control if users can increment/set their own stats
// - Admins to set user statistics
// - Users to increment/setup their own stats
// --- Nice to haves
func rpcmd_stat_init(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: stat_init")
	// Checks if feature is set to active
	// - look up serverDb to see if server exists/what feature toggles it has
	// - if server doesn't exist as a record, add it
	// - if server doesn't have feature enabled, ignore this
	// - if server has feature enabled, check that stats exist/set them if it doesn't
	// If set to active, checks if system has stats initiated

	// If not, goes ahead and activates it
}

func rpcmd_stat_up(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: statUp")

}

// ********** Titles System ********** //
// - Admins are able to set user titles
// - Titles are viewable within profiles
// --- Nice to Haves
// - Admins are able to set up automatic roles when conditions are met

// ********** Utility Functions ********** //
// - util_<?feature>_<featurecommand>

func msg_rpcmd_help() string {
	message := "*Here is what you can do with the freeform inventory system*"
	message += "\n```Markdown"
	message += "\n# == Commands List == #"
	message += "\n== User account"
	message += "\no $rpinventory    - Check your inventory"
	message += "\n                  - alias: $rpinv"
	message += "\no $rpigive @<user> <params> - Create an item and give it to a user"
	message += "\n                  - alias: $rpig"
	message += "\n                  - params: property:value - property2:value2"
	message += "\n                  - NOTE: Name expected and semi-required"
	message += "\n                  - example: $rpig @friend name:A present - description:It is a small white box wrapped in blue ribbon - weight: 0.5kg"
	message += "\no $rpidiscard ID# - Remove an item from your own inventory"
	message += "\n                  - Please use the ID displayed next to the item in your inventory"
	message += "\n                  - alias: $rpid"
	message += "```"

	return message
}

func util_checkIfValidUserPing(userPing string) error {
	log.Println("= On checkIfValidUser")
	log.Printf("\nReceived: %s", userPing)
	r, err := regexp.Compile("<[@!]+[0-9]+>")
	if err != nil {
		log.Println("Regexp unsuccessfully created in checkIfValidUser")
		log.Printf("\nError: \n%+v", err.Error())
		return errors.New("Invalid regexp")
	}

	result := r.FindString(userPing)
	if result == "" {
		log.Println("A user was not given to checkIfValidUser")
		return errors.New("Not a user")
	} else {
		return nil
	}
}

func util_getGuildId(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	channelInfo, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println("Cannot get server info")
		log.Printf("Error: %+v", err.Error())
		return "", err
	}
	guildId := channelInfo.GuildID
	return guildId, nil
}

func util_getParams_itemGive(msg string, withPing bool) (map[string]string, error) {
	var params = make(map[string]string)
	var err error
	// Check which command, if it's a shortcut
	cmdString := ""
	if strings.HasPrefix(msg, CMD_PREFIX+"rpigive") {
		cmdString = CMD_PREFIX + "rpigive "
	} else if strings.HasPrefix(msg, CMD_PREFIX+"rpig ") {
		cmdString = CMD_PREFIX + "rpig "
	}

	log.Printf("Got this msg: %s", msg)
	msgSansCmd := msg[len(cmdString):]

	index := strings.Index(msgSansCmd, " ")
	giveTo := msgSansCmd[:index]
	itemProps := msgSansCmd[index+1:]

	if withPing == true {
		err = util_checkIfValidUserPing(giveTo)
		if err != nil {
			log.Println("Error: This isn't a user ping")
			log.Printf("%s", err.Error())
			params["user"] = strings.Trim(giveTo, "@")
			return nil, err
		}

		r, err := regexp.Compile("[0-9]+")
		if err != nil {
			log.Println("Regexp unsuccessful init")
			log.Printf("%s", err.Error())
			return nil, err
		}
		userId := r.Find([]byte(giveTo))
		if userId == nil {
			log.Println("userId not found in ping ?? somehow")
			return nil, err
		}
		params["user"] = string(userId)
	} else {
		params["user"] = giveTo
	}

	// Parsing item commands
	params["item"] = itemProps

	return params, nil
}
