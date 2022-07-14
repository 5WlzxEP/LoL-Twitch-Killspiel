package main

import (
	"Killspiel"
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	config := getConfig("config.json")
	f, err := os.OpenFile(config.Logpath+"killspiel.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening log-file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	config.State = Killspiel.Idle
	Killspiel.SetConfig(config)

	//log.Printf("%v, %v, %s, %s, %v", channel, login, username, oath, err)

	config.TwitchClient = twitch.NewClient(config.Username, config.Oath)
	c := make(chan twitch.PrivateMessage)
	go Killspiel.Message(c)
	config.TwitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if config.State == Killspiel.Wettphase {
			c <- message
		}
		//log.Println(message.Message)
	})

	config.TwitchClient.Join(config.Twitchchannel)
	if config.Joinmessage {
		config.TwitchClient.Say(config.Twitchchannel, "Killspielbot aktiv")
	}

	//log.Println()
	go Killspiel.Statecontroll(Killspiel.GetLolID(config.Lolaccountname))

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
	defer f.Close()
	bites, _ := ioutil.ReadAll(f)
	conf := &Killspiel.Config{}
	err = json.Unmarshal(bites, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
