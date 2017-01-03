package emotes

import (
	// "github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

func LookUpThis(message string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	maxNum := len(EMOTE_LOOKUP)
	randEmote := EMOTE_LOOKUP[rand.Intn(maxNum-1)+1] // so numbers start from 1
	newMsg := randEmote + "\n" + message
	return newMsg
}
