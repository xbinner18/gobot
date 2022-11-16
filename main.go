package main

import (
	"log"
	"encoding/json"
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
        "github.com/xbinner18/gobot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
        type Response struct {
		Number struct {
		} `json:"number"`
		Scheme  string `json:"scheme"`
		Type    string `json:"type"`
		Brand   string `json:"brand"`
		Country struct {
			Numeric   string `json:"numeric"`
			Alpha2    string `json:"alpha2"`
			Name      string `json:"name"`
			Emoji     string `json:"emoji"`
			Currency  string `json:"currency"`
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
		} `json:"country"`
		Bank struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
		} `json:"bank"`
	}
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Panic("cannot load config")
	}
	token, channel := config.Token, config.ChannelID
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on Bot %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		re := regexp.MustCompile("[0-9]+")
		x := re.FindAllString(update.Message.Text, -1)
		match, _ := regexp.MatchString("[0-9]+", update.Message.Text)
		log.Println(match)
		
		if len(x) > 7 {
			continue
		}
		resp, err := http.Get("https://lookup.binlist.net/"+x[0][:6])
		if err != nil {
			log.Fatalln(err)
		}
                defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
                var response Response
                json.Unmarshal(body, &response)
		
		dump := tgbotapi.NewMessage(channel, fmt.Sprintf("%s\nBIN %s\n%s-%s-%s\n%s\n%s", update.Message.Text, x[0][:6], response.Scheme, response.Brand, response.Type, response.Bank.Name, response.Country.Name))

		if len(update.Message.Text) < 25 {
			continue
		}

		if match != false {
			time.Sleep(5 * time.Second) // time sleep 5 sec
			bot.Send(dump)
		}
	}
}
