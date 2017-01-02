package profile

import (
	// "encoding/json"
	"errors"
	// "github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
)

type UserData struct {
	Username string
	Credits  int
	Profile  string
	Title    string
}

func AddCredits(amount int) {

}

func CheckStats(username string) (string, error) {
	log.Printf("=== Checking status for: %s", username)
	resp, err := http.Get(API_ENDPOINT + "find/?name=" + username)
	if err != nil {
		log.Println("Unable to hit GET find/?name=")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		log.Fatal(err)
	}
	log.Printf("\n%s\n", body)

	if string(body) == ERR_NOUSER {
		log.Println("User doesn't exist. Calling register function.")
		newBody, err := RegisterUser(username)
		if err != nil {
			log.Println("Unable to register this user")
			log.Fatal(err)
		}
		log.Println("Finished calling register function.")
		return newBody, nil
	}

	// jsonResp := json.Unmarshal(resp, UserData)

	return string(body), nil
}

func RegisterUser(username string) (string, error) {
	log.Printf("=== Registering this user: %s", username)
	resp, err := http.Get(API_ENDPOINT + "add/user/?name=" + username)
	if err != nil {
		log.Println("Unable to hit GET add/?name=")
		log.Fatal(err)
		return "", errors.New("Unable to hit GET add/?name=")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		log.Fatal(err)
		return "", errors.New("Unable to read response body")

	}
	log.Printf("\n%s\n", body)

	return string(body), nil
}
