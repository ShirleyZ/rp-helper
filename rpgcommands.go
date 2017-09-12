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

// ********** Structured Inventory System ********** //
// - Structured items - pre-defined, ie iron ore
// - Utility func to give users items
// - Utility func to remove user items
// - Users can check their own inventory
// --- Nice to haves
// - Users can set showcase items viewable by other users

// ********** Job System ********** //
// --- Nice to haves
// - Users can choose/set their own profession path
// - Users can level up through profession path as they hit required levels/quests
// - Implement profession paths with their:
//		> collectable resources
// 		> craftable items

// ********** Gathering System ********** //
// --- RPG mechanic
// - Users are able to randomly collect resource of their profession while chatting

// ********** Market/Economy System ********** //
// --- Nice to haves
// - Users are able to sell/buy from the market
// - Determine how markets will work
//

// Considerations
// - Need to consider motivation to continue/goals to strive towards
