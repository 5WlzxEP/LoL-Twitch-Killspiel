package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"math/rand"
	"os"
	"time"
)

type conf struct {
	TwitchUser  string `json:"twitchUser"`
	TwitchOAuth string `json:"twitchOAuth"`
	Channel     string `json:"channel"`
}

func main() {

	rand.Seed(time.Now().UnixNano())

	f, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	data, _ := io.ReadAll(f)
	var config conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ready", config.TwitchUser)
	client := twitch.NewClient(config.TwitchUser, config.TwitchOAuth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil." {
			client.Say(config.Channel, fmt.Sprintf("!vote %d", rand.Intn(15)))
		}
	})
	client.Join(config.Channel)
	err = client.Connect()
	if err != nil {
		fmt.Println(err)
	}
}
