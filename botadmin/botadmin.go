package botadmin

import (
	// "github.com/bwmarrin/discordgo"
	// "strings"
	"../errors"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const API_ENDPOINT = "http://localhost:8080/api/"

var features = [4]string{"profile", "rpinv", "roll", "rpg"}

type BotConfig struct {
	GuildSettings map[string]GuildSettings
}

type GuildSettings struct {
	GuildId         string
	FeatureSettings map[string]map[string]string
}

func InitBotConfig() BotConfig {
	botConf := BotConfig{}
	botConf.GuildSettings = make(map[string]GuildSettings)
	// _ := getAllGuildConfig()
	guildSettingList := getAllGuildConfig()
	if len(guildSettingList) > 0 {
		botConf.GuildSettings = guildSettingList
	}
	return botConf
}

func InitGuildConfig(botConf *BotConfig, guildId string) {
	log.Printf("= botadmin:InitGuildConfig: GuildId %s", guildId)
	newConfig := GuildSettings{guildId, make(map[string]map[string]string)}
	log.Println("Setting stuff")
	for _, value := range features {
		log.Println(value)
		log.Printf("%+v", newConfig.FeatureSettings)
		newConfig.FeatureSettings[value] = make(map[string]string)
		newConfig.FeatureSettings[value]["enabled"] = "true"
	}
	log.Println("Setting specific settings")
	// Default profile settings
	newConfig.FeatureSettings["profile"]["creditOneAlias"] = "Splots"
	newConfig.FeatureSettings["profile"]["creditTwoAlias"] = "Ink Pots"
	newConfig.FeatureSettings["profile"]["creditThreeAlias"] = "Paintbrushes"

	botConf.GuildSettings[guildId] = newConfig
}

func FindSettingsByGuildId(guildId string, botConf *BotConfig) (GuildSettings, error) {
	log.Println("= botadmin:FindSettingsByGuildId: GuildId %s", guildId)
	if len(botConf.GuildSettings) <= 0 {
		return GuildSettings{}, errors.New(rperrors.ERR_NOGUILDSETTINGS)
	}

	guild := botConf.GuildSettings[guildId]
	log.Println(guild)
	if guild.GuildId == "" {
		InitGuildConfig(botConf, guildId)
		saveGuildConfig(botConf.GuildSettings[guildId])
	}

	return guild, nil
}

func GetConfigForGuild(guildId string, field string, botConf *BotConfig) (GuildSettings, error) {
	guild, err := FindSettingsByGuildId(guildId, botConf)
	if err != nil {
		return GuildSettings{}, errors.New("Couldn't find guild settings")
	}

	return guild, nil
}

func saveGuildConfig(guild GuildSettings) {
	log.Println("= botadmin:saveGuildConfig")
	// Make body
	// send req
	guildMarshalled, err := json.Marshal(guild)
	sendBody, err := url.ParseQuery(string(guildMarshalled))
	log.Printf("This is the body being sent to save guild\n%+v", sendBody)
	if err != nil {
		log.Println("Cannot parse query")
		log.Fatal(err)
	}

	_, err = http.PostForm(API_ENDPOINT+"botadmin/saveGuildConfig", sendBody)
	if err != nil {
		log.Println("didn't scucess the post")
		log.Fatal(err)
	}

	// Make req to the endpoint with body
	// resp, err := http.Get(API_ENDPOINT + "botadmin/initBotConfig")
	// if err != nil {
	//  log.Println("Unable to hit GET find/?name=")
	//  log.Fatal(err)
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//  log.Println("Unable to read response body")
	//  log.Fatal(err)
	// }
	// log.Printf("%s", body)
}

func getAllGuildConfig() map[string]GuildSettings {
	log.Println("= botadmin:getAllGuildConfig")
	resp, err := http.Get(API_ENDPOINT + "botadmin/getAllGuildConfig")
	if err != nil {
		log.Println("Unable to hit GET find/?name=")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		log.Fatal(err)
	}
	log.Println("Logging the response")
	log.Printf("\n\n%s\n\n", body)
	return make(map[string]GuildSettings)
}

// func SetConfig(botConf BotConfig, guildId string, field string, setTo bool) error {
//  log.Println("= botadmin:SetConfig: Setting %s to %t for guild %s", field, setTo, guildId)
//  // Finds the guild setting object
//  guildSettings, err := findSettingByGuildId(guildId, botConf)
//  if err != nil {
//     if (err == )
//    log.Printf("%v", err)
//    return errors.New("Could not retrieve guild settings: %v", err)
//  }

//  // Sets the field to the value
// }
