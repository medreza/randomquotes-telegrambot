package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quotescsvparser"
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
	records, err := quotescsvparser.ReadQuotesCsvFile("quotes.csv")
	if err != nil {
		log.Fatal(err)
	}
	quoteAuthor, quoteText := quotescsvparser.GetRandomQuote(records)
	quote := "<i>" + quoteText + "</i>" + "  -" + quoteAuthor
	sendMessageAPI := url + "/sendMessage?chat_id=" + chatID + "&text=" + quote + "&parse_mode=HTML"
	return getJSON(sendMessageAPI, target)
}

func periodicMessaging(apiURL string, t time.Time) error {
	userMsg := &getUpdates{}
	sentMsg := &sendMessage{}
	newMsgFlag := false
	lastBotUpdateID := 0

	getUpdatesAPI := apiURL + "/getUpdates"
	err := getJSON(getUpdatesAPI, &userMsg)
	if err != nil {
		return err
	}

	for _, value := range userMsg.Result {
		if value.UpdateID > 0 {
			lastBotUpdateID = value.UpdateID
			fmt.Print(t.Format(time.RFC850) + ": ")
			fmt.Print("Message received! UpdateID: " + strconv.Itoa(lastBotUpdateID))
			fmt.Print(" chatID: " + strconv.Itoa(value.Message.Chat.ID))
			fmt.Print(" from: " + value.Message.From.Username)
			newMsgFlag = true
			chatID := strconv.Itoa(value.Message.Chat.ID)
			err = sendMessageToUser(apiURL, chatID, &sentMsg)
			if err != nil {
				return err
			}
			fmt.Println(" - replied!")
		}
	}

	if !newMsgFlag {
		fmt.Print(t.Format(time.RFC850) + ": ")
		fmt.Println("No new messages...")
		return err
	}

	nextMsgUpdate := getUpdatesAPI + "?offset=" + strconv.Itoa(lastBotUpdateID+1)
	_, err = client.Get(nextMsgUpdate)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	apiURL := "https://api.telegram.org/bot" + botToken
	botIdentity := &getMe{}
	getMeAPI := apiURL + "/getMe"
	err = getJSON(getMeAPI, &botIdentity)
	if err != nil {
		return
	} else if botIdentity.Result.IsBot {
		fmt.Println("Bot name: " + botIdentity.Result.FirstName)
		fmt.Println("Bot username: " + botIdentity.Result.Username)
	} else {
		log.Fatal("Provided token is not a Telegram bot token.")
	}

	// Check message update every 5 seconds
	for t := range time.NewTicker(5 * time.Second).C {
		err = periodicMessaging(apiURL, t)
		if err != nil {
			fmt.Println()
			fmt.Println("Error: " + err.Error())
		}
	}

}
