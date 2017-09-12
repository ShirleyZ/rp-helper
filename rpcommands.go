package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

	params := util_getParams_itemGive(m.Content)
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
// - util_<feature>_<featurecommand>
func util_getParams_itemGive(msg string) []byte {

}
