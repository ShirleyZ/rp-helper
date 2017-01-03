package profile

import (
	// "github.com/bwmarrin/discordgo"
	// "encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type UserData struct {
	Username string
	Credits  int
	Profile  string
	Title    string
}

func AddCredits(username string, amount int) error {
	log.Printf("=== Adding %d credits to %s", amount, username)
	sendBody, err := url.ParseQuery("username=" + username + "&amount=" + strconv.Itoa(amount))
	if err != nil {
		log.Println("\nwewedwe\n")
		log.Fatal(err)
	}
	resp, err := http.PostForm(API_ENDPOINT+"credits/add/", sendBody)
	if err != nil {
		log.Println("didn't scucess the post")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("didn't suceess reading the body")
		log.Fatal(err)
	}
	if body == nil {
		return nil
	} else {
		return errors.New("Unsuccessful")
	}
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
		log.Println("User doesn't exist.")
		return "", errors.New(ERR_NOUSER)
		// newBody, err := RegisterUser(username)
		// if err != nil {
		// 	log.Println("Unable to register this user")
		// 	log.Fatal(err)
		// }
		// log.Println("Finished calling register function.")
		// return newBody, nil
	}

	// jsonResp := json.Unmarshal(resp, UserData)

	return string(body), nil
}

func RegisterUser(username string) (string, error) {
	log.Printf("=== Registering this user: %s", username)
	// Check that it doesn't currently exist
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
