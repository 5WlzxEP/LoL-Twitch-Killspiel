package main

import (
	"encoding/json"
	"fmt"
	"github.com/5WlzxEP/LoL-Twitch-Killspiel/src/Killspiel"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"log"
	"os"
)

func main() {
	config := getConfig("config.json")
	f, err := os.OpenFile(config.Logpath+"killspiel.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening log-file: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error occurded closing log file: %v", err)
		}
	}(f)
	log.SetOutput(f)
	log.Println("Starting...")

	config.State = Killspiel.Idle
	Killspiel.SetConfig(config)

	// Init Twitch Client
	config.TwitchClient = twitch.NewClient(config.Username, config.Oath)
	c := make(chan twitch.PrivateMessage)
	go Killspiel.Message(c)
	config.TwitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if config.State == Killspiel.Wettphase {
			c <- message
		}
	})

	config.TwitchClient.Join(config.Twitchchannel)
	if config.Joinmessage {
		log.Println("Sending Joinmessage")
		config.TwitchClient.Say(config.Twitchchannel, "Killspielbot aktiv")
	}

	//log.Println()
	go Killspiel.StateControl(Killspiel.GetLolID(config.Lolaccountname))

	err = config.TwitchClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func getConfig(file string) *Killspiel.Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error occurded closing config: %v\n", err)
		}
	}(f)
	bites, _ := io.ReadAll(f)
	conf := &Killspiel.Config{}
	err = json.Unmarshal(bites, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
