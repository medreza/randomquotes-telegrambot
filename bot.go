package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Bot information
type getMe struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"result"`
}

// Bot message update
type getUpdates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type sendMessage struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

// Define yout Telegram Bot token
var botToken = "<Your bot token>"

var client = &http.Client{Timeout: 5 * time.Second}

func getJSON(url string, target interface{}) error {
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

func sendMessageToUser(url string, chatID string, target interface{}) error {
	quoteText := "<i>Drangan trakioiot HAHAHA</i>   -Isyana Sarasvati"
	sendMessageAPI := url + "/sendMessage?chat_id=" + chatID + "&text=" + quoteText + "&parse_mode=HTML"
	return getJSON(sendMessageAPI, target)
}

func periodicMessaging(apiURL string, err error) {
	userMsg := &getUpdates{}
	sentMsg := &sendMessage{}
	newMsgFlag := false
	lastBotUpdateID := 0

	getUpdatesAPI := apiURL + "/getUpdates"
	err = getJSON(getUpdatesAPI, &userMsg)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, value := range userMsg.Result {
		if value.UpdateID > 0 {
			lastBotUpdateID = value.UpdateID
			fmt.Println("Message received!")
			fmt.Println(lastBotUpdateID)
			newMsgFlag = true
			chatID := strconv.Itoa(userMsg.Result[0].Message.Chat.ID)
			err = sendMessageToUser(apiURL, chatID, &sentMsg)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}

	if !newMsgFlag {
		fmt.Println("No new messages!")
		return
	}

	getNextUpdate := getUpdatesAPI + "?offset=" + strconv.Itoa(lastBotUpdateID+1)
	_, err = client.Get(getNextUpdate)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	var err error
	apiURL := "https://api.telegram.org/bot" + botToken
	botIdentity := &getMe{}
	getMeAPI := apiURL + "/getMe"
	err = getJSON(getMeAPI, &botIdentity)
	if err != nil {
		log.Fatal(err)
		return
	} else if botIdentity.Result.IsBot {
		fmt.Println("Bot name: " + botIdentity.Result.FirstName)
		fmt.Println("Bot username: " + botIdentity.Result.Username)
	} else {
		fmt.Println("Provided token is not a Telegram bot token")
		return
	}

	// Check message update every 5 seconds
	for t := range time.NewTicker(5 * time.Second).C {
		periodicMessaging(apiURL, err)
		fmt.Println(t)
	}

}