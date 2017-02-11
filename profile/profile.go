package profile

import (
	// "github.com/bwmarrin/discordgo"
	// "encoding/json"
	// "../errors"
	"encoding/json"
	"errors"
	"github.com/gorilla/Schema"
	// "html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const COOKIE_COST = 20

type UserData struct {
	Cookies  int
	Id       string
	Username string
	Credits  int
	Profile  string
	Title    string
}

func AddCredits(id string, username string, amount int) error {
	username = strings.ToLower(username)
	log.Printf("=== Adding %d credits to %s (%s)", amount, username, id)

	sendBody, err := url.ParseQuery("username=" + username + "&amount=" + strconv.Itoa(amount) + "&id=" + id)
	if err != nil {
		log.Println("Cannot parse query")
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

func CheckStats(id string) (string, error) {
	// username = strings.ToLower(username)
	log.Printf("=== Checking status for: %s", id)

	resp, err := http.Get(API_ENDPOINT + "find/?id=" + id)
	if err != nil {
		log.Println("Unable to hit GET find/?name=")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response body")
		log.Fatal(err)
	}
	log.Printf("%s", body)

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

func GiveCookie(giverId string, recipientId string, amount int) error {
	log.Println("== Giving Cookie ==")
	log.Printf("From: (%s) To: (%s) Amount: %d", giverId, recipientId, amount)
	if amount < 0 {
		return errors.New("Amount below 0")
	}
	// Get giver info
	data, err := CheckStats(giverId)
	giverInfo := UserData{}
	err = json.Unmarshal([]byte(data), &giverInfo)
	if err != nil {
		log.Printf("\n%v\n", err)
	}
	log.Printf("GiverInfo: %+v", giverInfo)

	// deduct amount from giver
	if giverInfo.Credits < COOKIE_COST*amount {
		return errors.New("Not enough credits")
	} else {
		giverInfo.Credits -= COOKIE_COST * amount

		// Update giver info
		sendInfo := url.Values{}
		encoder := schema.NewEncoder()
		err = encoder.Encode(giverInfo, sendInfo)
		_, err = http.PostForm(API_ENDPOINT+"profile/update/", sendInfo)
		if err != nil {
			return errors.New("Unsuccessful credit deduct")
		}

		// Get recipient info
		data, err := CheckStats(recipientId)
		recipInfo := UserData{}
		err = json.Unmarshal([]byte(data), &recipInfo)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		log.Printf("RecipInfo: %+v", recipInfo)
		recipInfo.Cookies += amount

		// Update recipient info
		sendInfo = url.Values{}
		err = encoder.Encode(recipInfo, sendInfo)
		_, err = http.PostForm(API_ENDPOINT+"profile/update/", sendInfo)
		if err != nil {
			return errors.New("Unsuccessful cookie give")
		}
		return nil
	}

}

func RegisterUser(id string, username string) (string, error) {
	username = strings.ToLower(username)
	log.Printf("== Registering this user: %s (%s)", username, id)
	// Check that it doesn't currently exist
	log.Println("Looking for user in database...")
	resp, err := http.Get(API_ENDPOINT + "find/?name=" + username + "&id=" + id)
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
	resp, err = http.Get(API_ENDPOINT + "add/user/?name=" + username + "&id=" + id)
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

func SetProfile(id string, profileBody string) (string, error) {
	log.Printf("== Registering this user: %s ", id)

	profileBody = url.QueryEscape(profileBody)
	sendBody, err := url.ParseQuery("id=" + id + "&profile=" + profileBody)
	if err != nil {
		log.Println("Parse query error")
		log.Fatal(err)
	}
	resp, err := http.PostForm(API_ENDPOINT+"profile/edit/", sendBody)
	if err != nil {
		log.Println("didn't scucess the post")
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("didn't suceess reading the body")
		log.Fatal(err)
	}
	log.Printf("Body: %+v", body)
	if body == nil {
		return "", nil
	} else {
		return "", errors.New("Unsuccessful")
	}
}
