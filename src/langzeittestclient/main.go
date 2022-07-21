package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"math/rand"
	"os"
)

type conf struct {
	TwitchUser  string `json:"twitchUser"`
	TwitchOAuth string `json:"twitchOAuth"`
	Channel     string `json:"channel"`
}

func main() {

	f, _ := os.Open("config.json")
	data, _ := io.ReadAll(f)
	var config conf
	_ = json.Unmarshal(data, &config)

	client := twitch.NewClient(config.TwitchUser, config.TwitchOAuth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil." {
			client.Say(config.Channel, fmt.Sprintf("!vote %d", rand.Intn(15)))
		}
	})
}
