package profile

import (
	// "github.com/bwmarrin/discordgo"
	// "encoding/json"
	// "../errors"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UserData struct {
	Username string
	Credits  int
	Profile  string
	Title    string
}

func AddCredits(username string, amount int) error {
	username = strings.ToLower(username)
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
	username = strings.ToLower(username)
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
	}
	return string(body), nil
}

// func CheckUserExists(username string) bool {
// 	username = strings.ToLower(username)
// 	log.Printf("=== Checking for this user: %s", username)

// 	resp, err := http.Get(API_ENDPOINT + "find/?name=" + username)
// }

func RegisterUser(username string) (string, error) {
	username = strings.ToLower(username)
	log.Printf("== Registering this user: %s", username)
	// Check that it doesn't currently exist
	log.Println("Looking for user in database...")
	resp, err := http.Get(API_ENDPOINT + "find/?name=" + username)
	log.Println("Done")

	log.Println("Reading GET response...")
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Println("Unable to read response body")
		log.Fatal(err2)
		return "", err2
	}
	log.Println("Done")

	if string(body) == ERR_USEREXISTS {
		// Shouldn't find any existing users with same name
		log.Println("User already exists")
		return "", errors.New(ERR_USEREXISTS)

	} else if err != nil {
		log.Println("Other error from find endpoint")
		log.Printf("\nError\n%v\n", err)
		return "", err
	}

	// Going ahead w/ account creation
	log.Println("Calling registration endpoint...")
	resp, err = http.Get(API_ENDPOINT + "add/user/?name=" + username)
	if err != nil {
		log.Println("Unable to hit GET add/?name=")
		log.Fatal(err)
		return "", errors.New("Unable to hit GET add/?name=")
	}
	log.Println("Done")

	log.Println("Reading response body...")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		log.Fatal(err)
		return "", errors.New("Unable to read response body")

	}
	log.Println("Done")
	log.Printf("\nNew user data:\n%s\n", body)

	return string(body), nil
	// }

}
