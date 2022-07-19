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
	config, fehlers, krit := getConfig("config.json")
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

	for _, fehler := range *fehlers {
		log.Println(fehler)
	}
	if krit {
		return
	}

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

func getConfig(file string) (*Killspiel.Config, *[]string, bool) {
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
	var confFehler []string
	var confFehlerKrit bool
	if conf.Logpath == "" {
		log.Fatalln("Logpath can't be empty")
	}
	if conf.Oath == "" {
		//log.Fatalln("OAuth can't be empty")
		confFehler = append(confFehler, "OAuth can't be empty")
		confFehlerKrit = true
	}
	if conf.Username == "" {
		//log.Fatalln("Username can't be empty")
		confFehler = append(confFehler, "Username can't be empty")
		confFehlerKrit = true
	}
	if conf.Twitchchannel == "" {
		//log.Fatalln("Kein Twitchchannel gesetzt, Nutzung des Programms ohne diesen Sinnlos.")
		confFehler = append(confFehler, "Kein Twitchchannel gesetzt, Nutzung des Programms ohne diesen Sinnlos.")
		confFehlerKrit = true
	}
	if conf.Lolaccountname == "" {
		//log.Fatalln("Kein LolAccount gesetzt => kein Sinn der Software.")
		confFehler = append(confFehler, "Kein LolAccount gesetzt => kein Sinn der Software.")
		confFehlerKrit = true
	}
	if conf.Lolapikey == "" {
		//log.Fatalln("Kein League API Key gesetzt, keine verfolgung möglich")
		confFehler = append(confFehler, "Kein League API Key gesetzt, keine Verfolgung möglich.")
		confFehlerKrit = true
	}

	if conf.Wettdauer == 0 {
		conf.Wettdauer = 120
		log.Println("Wettdauer nicht gesetzt, wurde auf 120s gesetzt.")
		confFehler = append(confFehler, "OAuth can't be empty")
	}

	return conf, &confFehler, confFehlerKrit
}
