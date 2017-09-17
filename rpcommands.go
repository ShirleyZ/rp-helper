package main

import (
	// "encoding/json"
	"fmt"
	"log"
	// "math/rand"
	"regexp"
	// "strconv"
	// "io/ioutil"
	"errors"
	"net/http"
	"net/url"
	"strings"
	// "time"

	"github.com/bwmarrin/discordgo"
)

const API_ENDPOINT = "http://localhost:8080/api/"

// ********** Bot Administration System ********** //
// - Feature toggles
// - Settings for commands/systems

// ********** Freeform Inventory System ********** //
// - Admin is able to remove someone's item
// - User is able to give themselves a new item
// - User is able to give someone else a new item
// - User is able to transfer their own item to someone else
// - User is able to remove their own items
// - Items are to be displayed on user's profile
// - User can see their own inventory
// --- Nice to haves
// - Items to keep history of ownership
// - Users can see other user's inventory
// - Users can hide items in their inventory (in a bag?)
// - Admin is able to remove someone's item with a reason
// - Admin logging functionality
// --- A/C
// - Items need to have a name, description, date created field
func rpcmd_item_check(s *discordgo.Session, m *discordgo.MessageCreate) {
	// User gets to check their inventory
}

func rpcmd_item_discard(s *discordgo.Session, m *discordgo.MessageCreate) {
	// User removes their own item
}

func rpcmd_item_give(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("=== Executing rpcmd: giveItem")

	params, err := util_getParams_itemGive(m.Content)
	if err != nil {
		log.Println("Error: problem getting params for item give")
		// TODO: handle error
	}

	fmt.Printf("\nparams:\n%+v\n", params)
	// TODO: bot-end checking of param

	sendBody, err := url.ParseQuery("userid=" + params["user"] + "&itemparams=" + params["item"])
	if err != nil {
		log.Println("Cannot parse query")
		log.Fatal(err)
	}
	url := API_ENDPOINT + "rpcmd/item/give/"
	fmt.Printf("\nHitting endpoint: %s\n", url)
	_, err = http.PostForm(url, sendBody)
	// resp, err := http.PostForm(API_ENDPOINT+"/rpcmd/item/give/", sendBody)
	if err != nil {
		log.Println("didn't success the post")
		log.Printf("%+v", err)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("didn't suceess reading the body")
	// 	log.Fatal(err)
	// }
	// if body == nil {
	// 	return nil
	// } else {
	// 	return errors.New("Unsuccessful")
	// }
	// Check who it goes to

	// hit api to do stuff

	// - first argument forced to be recipient
	// !giveItem @zaz #name A book #desc A really bright and glary book because i'm a good pal #weight 1kg
	// !giveItem @zaz name:A book / desc: A really bright and glary book because i'm a good pal / weight: 1kg
	// !giveItem @zaz - name:A book - desc: A really bright and glary book because i'm a good pal - weight: 1kg
	// !giveItem @zaz name:A book - desc: A really bright and glary book because i'm a good pal - weight: 1kg
	// !giveItem @zaz, name:A book, desc: A really bright and glary book because i'm a good pal
	// !giveItem [@zaz] name:[A book] desc:[This is another book about smug assholes, I thought you might like it]

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
func util_checkIfValidUser(userPing string) error {
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

func util_getParams_itemGive(msg string) (map[string]string, error) {
	var params = make(map[string]string)
	var err error
	// Check which command, if it's a shortcut
	cmdString := ""
	if strings.HasPrefix(msg, CMD_PREFIX+"giveitem") {
		cmdString = CMD_PREFIX + "giveitem "
	} else if strings.HasPrefix(msg, CMD_PREFIX+"gi ") {
		cmdString = CMD_PREFIX + "gi "
	}

	msgSansCmd := msg[len(cmdString):]

	index := strings.Index(msgSansCmd, " ")
	giveTo := msgSansCmd[:index]
	itemProps := msgSansCmd[index+1:]

	params["user"] = giveTo
	err = util_checkIfValidUser(giveTo)
	if err != nil {
		log.Println("Error: User is not valid")
		return nil, err
	}

	// Parsing item commands
	params["item"] = itemProps

	return params, nil
}
