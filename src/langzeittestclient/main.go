package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

type conf struct {
	TwitchUser  []string `json:"twitchUser"`
	TwitchOAuth []string `json:"twitchOAuth"`
	Channel     string   `json:"channel"`
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

	if len(config.TwitchUser) != len(config.TwitchOAuth) {
		panic("length of TwitchUser must be length of TwitchOAuth")
	}
	var wg sync.WaitGroup
	for i := range config.TwitchUser {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			client(config.TwitchUser[i], config.TwitchOAuth[i], &config)
		}(i)
	}
	wg.Wait()
}

func client(username string, oauth string, config *conf) {
	client := twitch.NewClient(username, oauth)
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil." {
			client.Say(config.Channel, fmt.Sprintf("!vote %d", rand.Intn(15)))
		}
	})
	client.Join(config.Channel)
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
	}
}
