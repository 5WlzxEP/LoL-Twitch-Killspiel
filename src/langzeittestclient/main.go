package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"log"
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
		log.Println(err)
	}

	if len(config.TwitchUser) != len(config.TwitchOAuth) {
		panic("length of TwitchUser must be length of TwitchOAuth")
	}
	var clients []*twitch.Client
	for i := range config.TwitchUser {
		clients = append(clients, twitch.NewClient(config.TwitchUser[i], config.TwitchOAuth[i]))
		//client(config.TwitchUser[i], config.TwitchOAuth[i], &config)
	}
	clients[0].OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil." {
			for _, client := range clients {
				client.Say(config.Channel, fmt.Sprintf("!vote %d", rand.Intn(15)))
			}
		}
	})

	var wg sync.WaitGroup
	wg.Add(len(clients))
	for _, client := range clients {
		client.Join(config.Channel)
		go func(client *twitch.Client) {
			defer wg.Done()
			err := client.Connect()
			if err != nil {
				log.Println(err)
			}
		}(client)
	}
	wg.Wait()
}
