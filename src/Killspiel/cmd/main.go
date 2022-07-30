package main

import (
	"encoding/json"
	"fmt"
	"github.com/5WlzxEP/LoL-Twitch-Killspiel/src/Killspiel"
	"github.com/gempir/go-twitch-irc/v2"
	"io"
	"log"
	"os"
	"time"
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
	id, err := Killspiel.GetLolID(config.Lolaccountname)
	if err != nil {
		log.Fatal(err)
	}
	go Killspiel.StateControl(id)

	var timeErr []time.Time

	for {
		for i := 0; i < 3; i++ {
			err = config.TwitchClient.Connect()
			if err != nil {
				log.Println(err)
			}
		}
		timeErr = append(timeErr, time.Now())
		if len(timeErr) > 5 {
			if timeErr[0].Add(5 * time.Minute).Before(timeErr[5]) {
				log.Fatal("More than 5 connection error with twitch in 5 Minutes, stopping!")
			}
			timeErr = timeErr[1:]
		}
	}
}

func getConfig(file string) (*Killspiel.GlobalConfig, *[]string, bool) {
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
	conf := &Killspiel.GlobalConfig{}
	err = json.Unmarshal(bites, conf)
	if err != nil {
		log.Fatal(err)
	}
	var confFehler []string
	var confFehlerKrit bool

	if conf.Oath == "" {
		confFehler = append(confFehler, "OAuth can't be empty")
		confFehlerKrit = true
	}
	if conf.Username == "" {
		confFehler = append(confFehler, "Username can't be empty")
		confFehlerKrit = true
	}
	if conf.Twitchchannel == "" {
		confFehler = append(confFehler, "Kein Twitchchannel gesetzt, Nutzung des Programms ohne diesen Sinnlos.")
		confFehlerKrit = true
	}
	if conf.Lolaccountname == "" {
		confFehler = append(confFehler, "Kein LolAccount gesetzt => kein Sinn der Software.")
		confFehlerKrit = true
	}
	if conf.Lolapikey == "" {
		confFehler = append(confFehler, "Kein League API Key gesetzt, keine Verfolgung möglich.")
		confFehlerKrit = true
	}

	if conf.Wettdauer == 0 {
		conf.Wettdauer = 120
		confFehler = append(confFehler, "Keine Wettdauer gesetzt")
	}

	if conf.LoLRegion == "" {
		conf.LoLRegion = "euw1"
		conf.LolServer = Killspiel.Europe
	} else {
		region, server, notFound := Killspiel.LoLRegionToServer(conf.LoLRegion)
		if notFound {
			confFehler = append(confFehler, "Falsche Region gesetzt.")
			confFehlerKrit = true
		} else {
			conf.LoLRegion = region
			conf.LolServer = server
		}
	}

	return conf, &confFehler, confFehlerKrit
}
