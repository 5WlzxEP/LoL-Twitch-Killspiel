package Killspiel

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Gamestate int

const (
	Idle Gamestate = iota
	Wettphase
	Spielphase
	Auswertungsphase
)

type Config struct {
	Username       string `json:"Username"`
	Oath           string `json:"Oath"`
	Wettdauer      int    `json:"Wettdauer"`
	Twitchchannel  string `json:"Twitchchannel"`
	Lolaccountname string `json:"Lolaccountname"`
	Lolapikey      string `json:"Lolapikey"`
	Joinmessage    bool   `json:"Joinmessage"`
	State          Gamestate
	TwitchClient   *twitch.Client
}

//type result struct {
//	Kills      int `json:"Kills"`
//	Teilnehmer struct {
//		Gewinnern struct {
//			Name string `json:"Name"`
//			Tipp int    `json:"Tipp"`
//		} `json:"Gewinnern"`
//		Verlierer struct {
//			Name string `json:"Name"`
//			Tipp int    `json:"Tipp"`
//		} `json:"Verlierer"`
//	} `json:"Teilnehmer"`
//}

type result struct {
	Kills      int        `json:"Kills"`
	Teilnehmer Teilnehmer `json:"Teilnehmer"`
}
type Teilnehmer struct {
	Gewinner  []Teilnehmer2 `json:"Gewinner"`
	Verlierer []Teilnehmer2 `json:"Verlierer"`
}

type Teilnehmer2 struct {
	Name string `json:"Name"`
	Tipp int    `json:"Tipp"`
}

var bessereDaten map[int][]string
var config *Config
var daten map[string]int
var wettdauer time.Duration

func SetConfig(config2 *Config) {
	config = config2
	wettdauer = time.Duration(config.Wettdauer)
	bessereDaten = map[int][]string{}
	aktuellesGame = &game{}
}

func Message(messagechan chan twitch.PrivateMessage) {
	//var mess *twitch.PrivateMessage
	for true {
		message := <-messagechan
		if strings.HasPrefix(message.Message, "!vote") {
			_, value, found := strings.Cut(message.Message, " ")
			if found {
				v, err := strconv.Atoi(value)
				if err == nil {
					daten[message.User.DisplayName] = v
					//log.Println(message.User.DisplayName, v)
				}
			}
		}

	}
}

func StarteWette() {
	config.State = Wettphase
	config.TwitchClient.Say(config.Twitchchannel, "/announce  Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil.")
	daten = make(map[string]int)
	time.Sleep(wettdauer * time.Second)
	config.State = Spielphase
	config.TwitchClient.Say(config.Twitchchannel, fmt.Sprintf("/announce  Killspiel-Teilnahme beendet, es haben %d Personen teilgenommen.", len(daten)))

	for player, points := range daten {
		bessereDaten[points] = append(bessereDaten[points], player)
	}
}

func Auswertung() {
	config.State = Auswertungsphase
	killd := GetKills()
	//log.Println(killd)
	//log.Println(aktuellesGame)
	var ind int = 11
	player := lolidToPuuid()
	for i := range killd.Metadata.Participants {
		if killd.Metadata.Participants[i] == player {
			ind = i
			break
		}
	}
	if ind == 11 {
		log.Fatal("player not found in result")
	}

	if killd.Info.Participants[ind].TeamEarlySurrendered || killd.Info.Participants[(ind+5)%10].TeamEarlySurrendered {
		config.TwitchClient.Say(config.Twitchchannel, "/announce Killspiel wurde abgebrochen, da Remaked wurde.")
	} else {
		kills := killd.Info.Participants[ind].Kills
		gewinner := bessereDaten[kills]
		res := result{Kills: kills, Teilnehmer: Teilnehmer{Gewinner: make([]Teilnehmer2, len(gewinner)), Verlierer: make([]Teilnehmer2, 0)}}
		for i, g := range gewinner {
			res.Teilnehmer.Gewinner[i] = Teilnehmer2{Name: g, Tipp: kills}
		}
		for k, l := range bessereDaten {
			if k == kills {
				continue
			}
			for _, i := range l {
				res.Teilnehmer.Verlierer = append(res.Teilnehmer.Verlierer, Teilnehmer2{Name: i, Tipp: k})
			}
		}

		fileInfo, err := os.Stat("results")
		if err != nil || !fileInfo.IsDir() {
			//log.Fatal(err)
			os.Mkdir("results", 0750)
		}

		file, err := os.Create(fmt.Sprintf("results/%d.json", aktuellesGame.matchId))
		if err != nil {
			log.Fatal(err)
		}
		bites, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		file.Write(bites)

		config.TwitchClient.Say(config.Twitchchannel,
			fmt.Sprintf("/announce Killspiel wurde beendet. %s hat %d gemacht. %s haben die richtige Killanzahl getippt.",
				config.Lolaccountname, kills, strings.Join(gewinner, ", ")))
	}
	daten = make(map[string]int)
	bessereDaten = map[int][]string{}
	config.State = Idle
}
