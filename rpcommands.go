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

	// Do a lookup of their own stuff
	endpoint := API_ENDPOINT + "rpcmd/item/check/"
	sendBody, err := url.ParseQuery("userid=" + userId)
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
	for key, value := range userInventory.Items {
		currItem := value
		log.Printf("%s: %+v ", key, currItem)
		message += "# [" + currItem["itemid"] + "] " + currItem["name"] + " #\n"
		message += "- Description: " + currItem["desc"] + "\n"
		for propName, propValue := range currItem {
			if (propName != "itemid") && (propName != "name") && (propName != "desc") {
				message += "- " + propName + ": " + propValue + "\n"
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
	// User removes their own item
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
	noPingParams, err := util_getParams_itemGive(noMentions, false)
	if err != nil {
		log.Println("Error: problem getting params for item give")
		return
		// TODO: handle error
	}
	fmt.Printf("\nparams:\n%+v\n", params)
	fmt.Printf("\nnopingparams:\n%+v\n", noPingParams)
	// TODO: bot-end checking of param

	sendBody, err := url.ParseQuery("userid=" + params["user"] + "&itemparams=" + noPingParams["item"])
	if err != nil {
		log.Println("Cannot parse query")
		log.Printf("Error: %+v", err.Error())
		return
	}
	url := API_ENDPOINT + "rpcmd/item/give/"
	fmt.Printf("\nHitting endpoint: %s\n", url)
	_, err = http.PostForm(url, sendBody)
	if err != nil {
		log.Println("didn't success the post")
		log.Printf("%+v", err)
		return
	}
	//message the user success
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("Unable to create private channel")
		log.Printf("\n%v\n", err)
		return
	}
	sendToThis := channel.ID
	_, err = s.ChannelMessageSend(sendToThis, "Item given to "+noPingParams["user"])

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
func util_checkIfValidUserPing(userPing string) error {
	log.Println("= On checkIfValidUser")
	log.Printf("\nReceived: %s", userPing)
	r, err := regexp.Compile("<@![0-9]+>")
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

func util_getParams_itemGive(msg string, withPing bool) (map[string]string, error) {
	var params = make(map[string]string)
	var err error
	// Check which command, if it's a shortcut
	cmdString := ""
	if strings.HasPrefix(msg, CMD_PREFIX+"giveitem") {
		cmdString = CMD_PREFIX + "giveitem "
	} else if strings.HasPrefix(msg, CMD_PREFIX+"gi ") {
		cmdString = CMD_PREFIX + "gi "
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
